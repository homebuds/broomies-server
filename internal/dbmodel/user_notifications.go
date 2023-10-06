package model

import "github.com/google/uuid"

type UserNotification struct {
	NotificationID uuid.UUID    `gorm:"not null" json:"notificationId"`
	AccountID      uuid.UUID    `gorm:"not null" json:"accountId"`
	Seen           bool         `gorm:"not null;default:false" json:"seen"`
	Notification   Notification `json:"notification"`
	Account        Account
}
