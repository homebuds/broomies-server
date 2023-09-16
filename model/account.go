package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()"`
	FirstName   string    `gorm:"size:255"`
	LastName    string    `gorm:"size:255"`
	Email       string    `gorm:"size:255, unique"`
	HouseholdID uuid.UUID
	Household   Household
}
