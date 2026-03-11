package model

import "time"

type VerificationProductionDetail struct {
	ID               uint   `gorm:"primaryKey"`
	VerificationUUID string `gorm:"type:char(36);not null"`

	ProductID int
	Woolen    int
	Blended   int
	Others    int
	Total     int

	CreatedAt time.Time
}