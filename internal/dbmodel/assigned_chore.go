package model

import (
	"time"

	"github.com/google/uuid"
)

type AssignedChore struct {
	ID        uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	ChoreID   uuid.UUID `gorm:"not null" json:"choreId"`
	AccountID uuid.UUID `gorm:"not null" json:"accountId"`
	Date      time.Time `gorm:"type:date" json:"dueDate"`
	Chore     Chore     `json:"chore"`
	Account   Account   `json:"account"`
	Completed bool      `json:"completed"`
}
