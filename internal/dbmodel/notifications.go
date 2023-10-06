package model

import "github.com/google/uuid"

type Notification struct {
	ID                     uuid.UUID  `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	AccountID              uuid.UUID  `gorm:"not null" json:"accountId"`
	HouseholdID            uuid.UUID  `gorm:"not null" json:"householdId"`
	Action                 string     `json:"action"`
	ActorAccountId         uuid.UUID  `gorm:"not null" json:"actorAccountId"`
	ActorChoreID           *uuid.UUID `json:"assignedChoreId"`
	FinancialTransactionID *uuid.UUID `json:"financialTransactionId"`
}
