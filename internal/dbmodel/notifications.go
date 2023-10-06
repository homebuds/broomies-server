package model

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID                     uuid.UUID  `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	Action                 string     `json:"action"`
	ActorAccountID         uuid.UUID  `gorm:"not null" json:"actorAccountId"`
	ActorChoreID           *uuid.UUID `json:"assignedChoreId"`
	FinancialTransactionID *uuid.UUID `json:"financialTransactionId"`
	ActorAccount           Account    `json:"actorAccount"`
	ActorChore             Chore      `json:"actorChore"`
	CreatedAt              time.Time  `gorm:"default: now()" json:"createdAt"`
}
