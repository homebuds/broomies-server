package model

import "github.com/google/uuid"

// RoommateScore is the model for roommate scores.
type RoommateScore struct {
	ID          uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	AccountID   uuid.UUID `gorm:"not null" json:"accountId"`
	HouseholdID uuid.UUID `gorm:"not null" json:"householdId"`
	Points      uint      `gorm:"not null" json:"points"`
	Account     Account   `json:"account"`
	Household   Household `json:"household"`
}
