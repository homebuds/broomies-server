package model

import "github.com/google/uuid"

// ChoreReview is the model for chore reviews.
type ChoreCompletionReview struct {
	UserNotificationID uuid.UUID `gorm:"not null" json:"userNotificationId"`
	Review             string    `gorm:"not null" json:"review"`
	UserNotification
}
