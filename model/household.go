package model

import "github.com/google/uuid"

type Household struct {
	ID   uuid.UUID `gorm:"primaryKey; default:uuid_generate_v4()"`
	Name string
}
