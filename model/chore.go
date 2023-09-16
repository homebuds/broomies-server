package model

import "github.com/google/uuid"

type JSONB map[string]interface{}

type Chore struct {
	ID          uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()"`
	Name        string    `gorm:"size:255"`
	Description string    `gorm:"size:255"`
	Points      uint      `gorm:"size:255"`
	Repetition  JSONB     `gorm:"type:jsonb"`
}
