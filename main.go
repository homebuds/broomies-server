package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	model "github.com/homebuds/broomies-server/internal/dbmodel"
	"github.com/homebuds/broomies-server/service"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

var db service.DBService

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	dsn := os.Getenv("POSTGRES_DSN")
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
	r.GET("/api/top-roommates/:household_id", GetRoommatesWithTopScore)
	r.PATCH("/api/assigned-chore/:assigned_chore_id/complete", MarkAssignedChoreComplete)
	r.PUT("/api/financial-transaction", CreateFinancialTransaction)
	r.GET("/api/financial-transaction/list/:household_id/:account_id", GetFinancialTransactions)
	r.GET("/api/spend-information/household/:household_id/account/:account_id", GetSpendInformation)
	r.Run() // listen and serve on
}

type LoginRequest struct {
	Email string `json:"email"`
}

func GetRoommatesWithTopScore(c *gin.Context) {
	householdParam := c.Param("household_id")
	householdID, err := uuid.Parse(householdParam)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var topN int
	if topNParam := c.Query("top_n"); topNParam == "" {
		topN = 4
	} else {
		topN, err = strconv.Atoi(topNParam)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	accounts, err := db.GetAccountsWitMostPoints(householdID, topN)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, accounts)

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
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	assignedChore, err := db.GetAssignedChore(assignedChoreUUID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	currChoreDueDate := assignedChore.Date
	points := int(assignedChore.Chore.Points)
	if currChoreDueDate.Before(time.Now()) {
		points *= -1
	}
	err = db.UpsertRoommateScore(assignedChore.AccountID, assignedChore.Chore.HouseholdId, points)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get user with least chores
	acc, err := db.GetAccountsWithLeastChores(assignedChore.Chore.HouseholdId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.CreateAssignedChore(assignedChore.Chore.ID, acc[0].ID, currChoreDueDate.AddDate(0, 0, 7))
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, assignedChore)
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
	Icon        string                `json:"icon"`
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
		Icon:           req.Icon,
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
		diff := day - (int64(curr.Weekday()))
		if diff <= 0 {
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
			diff := day - int64(curr.Weekday())
			if diff <= 0 {
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

type CreateFinancialTransactionRequest struct {
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Name        string    `json:"name"`
	AccountID   uuid.UUID `json:"accountId"`
	HouseholdID uuid.UUID `json:"householdId"`
}

type TransactionWithOwed struct {
	model.FinancialTransaction
	Owed float64 `json:"owed"`
}

func GetFinancialTransactions(c *gin.Context) {
	householdId := c.Param("household_id")
	householdUuid, err := uuid.Parse(householdId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	accountIdParam := c.Param("account_id")
	accountId, err := uuid.Parse(accountIdParam)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	transactions, err := db.GetUnsettledFinancialTransactions(householdUuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	accounts, err := db.GetAccounts(householdUuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	transactionsWithOwed := make([]TransactionWithOwed, len(transactions))
	for i, tx := range transactions {
		if tx.AccountID == accountId {
			transactionsWithOwed[i] = TransactionWithOwed{
				FinancialTransaction: tx,
				Owed:                 -(tx.Amount / float64(len(accounts))) * float64(len(accounts)-1),
			}
		} else {
			transactionsWithOwed[i] = TransactionWithOwed{
				FinancialTransaction: tx,
				Owed:                 tx.Amount / float64(len(accounts)),
			}
		}
	}
	c.JSON(200, transactionsWithOwed)
}

func CreateFinancialTransaction(c *gin.Context) {
	var req CreateFinancialTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, err := db.CreateFinancialTransaction(req.Amount, req.Name, req.AccountID, req.HouseholdID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, tx)
}

type SpendInformation struct {
	TotalSpent           float64 `json:"totalSpent"`
	AvgSpent             float64 `json:"avgSpent"`
	AmountOwed           float64 `json:"amountOwed"`
	RoommatePointsAmount float64 `json:"roommatePointsAmount"`
}

func GetSpendInformation(c *gin.Context) {
	householdId := c.Param("household_id")
	targetAccountId := c.Param("account_id")
	householdUuid, err := uuid.Parse(householdId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	targetAccountUUID, err := uuid.Parse(targetAccountId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	unsettledTransactions, err := db.GetUnsettledFinancialTransactions(householdUuid)
	accounts, err := db.GetAccounts(householdUuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	totalSpent := 0.0
	spentPerPerson := make(map[uuid.UUID]float64)
	for _, tx := range unsettledTransactions {
		totalSpent += tx.Amount
		if _, ok := spentPerPerson[tx.AccountID]; !ok {
			spentPerPerson[tx.AccountID] = 0.0
		}
		spentPerPerson[tx.AccountID] += tx.Amount
	}
	var meanSpentPerPerson float64 = 0
	if len(accounts) > 0 {
		meanSpentPerPerson = totalSpent / float64(len(accounts))
	}
	roommatePointsPerPerson := make(map[uuid.UUID]float64)
	for _, account := range accounts {
		accountID := account.ID
		points, err := db.GetPointsForAccount(accountID)
		if err != nil {
			points = 500
		}
		roommatePointsPerPerson[accountID] = float64(points)
		spentPerPerson[accountID] -= meanSpentPerPerson
	}
	pointAmount := totalSpent * 0.1
	totalPoints := 0.0
	for _, points := range roommatePointsPerPerson {
		totalPoints += points
	}
	log.Printf("Total points: %f", totalPoints)
	pointSavings := 0.0

	spentPerPerson[targetAccountUUID] -= pointSavings
	spendInfo := SpendInformation{
		TotalSpent: totalSpent,
	}
	if len(spentPerPerson) > 0 {
		spendInfo.AvgSpent = meanSpentPerPerson
	}
	if _, ok := spentPerPerson[targetAccountUUID]; ok {
		spendInfo.AmountOwed = spentPerPerson[targetAccountUUID]
	}
	if _, ok := roommatePointsPerPerson[targetAccountUUID]; ok && totalPoints > 0 {
		spendInfo.RoommatePointsAmount = pointAmount * (roommatePointsPerPerson[targetAccountUUID] / totalPoints)
	}
	c.JSON(200, spendInfo)
}
