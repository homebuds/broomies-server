package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ChoreRepetition struct {
	Days []int64 `json:"days"`
}

type Chore struct {
	ID             uuid.UUID     `gorm:"primaryKey; default:uuid_generate_v4()" json:"id"`
	Name           string        `gorm:"size:255" json:"name"`
	Description    string        `gorm:"size:255" json:"description"`
	Points         uint          `gorm:"size:255" json:"points"`
	Icon           string        `json:"icon"`
	HouseholdId    uuid.UUID     `gorm:"not null" json:"householdId"`
	WeekDayRepeats pq.Int64Array `gorm:"type:integer[]" json:"weekDayRepeats"`
	Household      Household     `json:"household"`
}
