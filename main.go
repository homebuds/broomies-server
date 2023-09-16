package main

import (
	"homebuds/model"
	"homebuds/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var db service.DBService

func main() {
	dsn := "postgresql://rustam:1DOh1353k5zPYL5-7HtecQ@htnx-project-12229.7tt.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	var err error
	db, err = service.NewDBService(dsn)
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.POST("/login", Login)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/api/chore/list/:household_id", GetChores)
	r.POST("/api/chore", CreateChore)
	r.GET("/api/assigned-chore/list/:household_id", GetAssignedChores)
	r.GET("/api/account/:account_id", GetAccountById)
	r.PATCH("/api/assigned-chore/:assigned_chore_id/complete", MarkAssignedChoreComplete)
	r.Run() // listen and serve on
}

type LoginRequest struct {
	Email string `json:"email"`
}

func GetAccountById(c *gin.Context) {
	accountId := c.Param("account_id")
	accountUuid, err := uuid.Parse(accountId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	account, err := db.GetAccountByID(accountUuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, account)
}

func MarkAssignedChoreComplete(c *gin.Context) {
	assignedChoreParam := c.Param("assigned_chore_id")
	assignedChoreUUID, err := uuid.Parse(assignedChoreParam)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = db.MarkAssignedChoreComplete(assignedChoreUUID)
	assignedChore, err := db.GetAssignedChore(assignedChoreUUID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	currChoreDueDate := assignedChore.Date
	points := int(assignedChore.Chore.Points)
	if currChoreDueDate.After(time.Now()) {
		points *= -1
	}

	// Get user with least chores
	_, err = db.GetAccountsWithLeastChores(assignedChore.Chore.HouseholdId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

}
func GetChores(c *gin.Context) {
	householdId := c.Param("household_id")
	householdUuid, err := uuid.Parse(householdId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	chores, err := db.GetChores(householdUuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, chores)
}
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	account, err := db.GetAccountByEmail(req.Email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, account)
}

type CreateChoreRequest struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Points      uint                  `json:"points"`
	HouseholdId uuid.UUID             `json:"householdId"`
	Repetition  model.ChoreRepetition `json:"repetition"`
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func GetAssignedChores(c *gin.Context) {
	householdId := c.Param("household_id")
	householdUuid, err := uuid.Parse(householdId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	chores, err := db.GetAssignedChores(householdUuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, chores)
}
func CreateChore(c *gin.Context) {
	var req CreateChoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(req.Repetition.Days) == 0 {
		c.JSON(400, gin.H{
			"error": "Repetition days should not be empty",
		})
		return
	}
	chore := model.Chore{
		Name:           req.Name,
		Description:    req.Description,
		Points:         req.Points,
		WeekDayRepeats: pq.Int64Array(req.Repetition.Days),
		HouseholdId:    req.HouseholdId,
	}
	err := db.CreateChore(&chore)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	curr := time.Now()
	var nextDate int64 = 7
	for _, day := range req.Repetition.Days {
		diff := day - (int64(curr.Weekday()) + 1)
		if diff < 0 {
			diff += 7
		}
		nextDate = min(diff, nextDate)
	}
	nextTaskOccurrence := curr.AddDate(0, 0, int(nextDate))
	// Get user with least chores
	acc, err := db.GetAccountsWithLeastChores(req.HouseholdId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Create assigned chore
	if len(chore.WeekDayRepeats) == 1 {
		err = db.CreateAssignedChore(chore.ID, acc[0].ID, nextTaskOccurrence)
	} else {
		nextTaskOccurrences := make([]time.Time, 0)
		for _, day := range req.Repetition.Days {
			diff := day - (int64(curr.Weekday()) + 1)
			if diff < 0 {
				diff += 7
			}
			nextTaskOccurrences = append(nextTaskOccurrences, curr.AddDate(0, 0, int(diff)))
		}
		for i := range nextTaskOccurrences {
			err = db.CreateAssignedChore(chore.ID, acc[i%len(acc)].ID, nextTaskOccurrences[i])
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
		}
	}
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, chore)
}
