package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignedChore struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()"`
	ChoreID   uuid.UUID
	AccountID uuid.UUID
	Date      time.Time `gorm:"type:date"`
	Chore     Chore
	Account   Account
}
