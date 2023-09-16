package main

import (
	"housebuds/model"
	"housebuds/service"
	"time"

	"github.com/gin-gonic/gin"
)

var db service.DBService

func main() {
	dsn := ""
	var err error
	db, err = service.NewDBService(dsn)
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.POST("/login", Login)

	r.Run() // listen and serve on
}

type GetAccountRequest struct {
	Email string `json:"email"`
}

func Login(c *gin.Context) {
	var req GetAccountRequest
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
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Points      uint                   `json:"points"`
	Repetition  map[string]interface{} `json:"repetition"`
}

func CreateChore(c *gin.Context) {
	var req CreateChoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	chore := model.Chore{
		Name:        req.Name,
		Description: req.Description,
		Points:      req.Points,
		Repetition:  req.Repetition,
	}
	err := db.CreateChore(&chore)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Load already assigned chores in the household

	// Assign chore to users based on repetition

	//
	// Create a monthly schedule for a chore, based of repetitions assign to users

	c.JSON(200, chore)
}

// Given an array of users, array of chore instances that occur on a given day of the week
// and whether task task occurs weekly or biweekly or monthly
// assign chore instances to users
func AssignChores(users []model.Account, chores []model.Chore) {
	now := time.Now()
	startDateOfWeek := int(now.Weekday())
	daysInMonth := DaysInMonth(now)
	weeks := make([][]int, 5)
	// Create array of weeks
	dayCount := 0
	for i := 0; i < 5; i++ {
		for (startDateOfWeek < 7) && (dayCount < daysInMonth) {
			weeks[i] = append(weeks[i], startDateOfWeek+1)
			startDateOfWeek++
			dayCount++
		}
		startDateOfWeek = 0
	}

}

func DaysInMonth(t time.Time) int {
	y, m, _ := t.Date()
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
