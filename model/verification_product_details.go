package model

import "time"

type VerificationProductDetail struct {
	ID                   uint   `gorm:"primaryKey"`
	VerificationUUID     string `gorm:"type:char(36);not null"`
	ProductSubcategoryID int
	Length               float64
	Width                float64

	CreatedAt time.Time
}