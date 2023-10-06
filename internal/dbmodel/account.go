package model

import (
	"github.com/google/uuid"
)

type Account struct {
	ID          uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	FirstName   string    `gorm:"size:255" json:"firstName"`
	LastName    string    `gorm:"size:255" json:"lastName"`
	Email       string    `gorm:"size:255, unique" json:"email"`
	HouseholdID uuid.UUID `gorm:"not null" json:"householdId"`
	PictureURL  string    `gorm:"size:255" json:"pictureUrl"`
	Household   Household
}
