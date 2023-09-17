package model

import (
	"time"

	"github.com/google/uuid"
)

type FinancialTransaction struct {
	ID          uuid.UUID  `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`
	AccountID   uuid.UUID  `json:"accountId"`
	HouseholdID uuid.UUID  `json:"householdId"`
	Account     Account    `json:"account"`
	Household   Household  `json:"household"`
	SettledAt   *time.Time `json:"settledAt" `
	CreatedAt   time.Time  `gorm:"default: now()" json:"createdAt"`
}
