package model

import "time"

type VerificationProductionCenter struct {
	ID               uint   `gorm:"primaryKey"`
	VerificationUUID string `gorm:"type:char(36);not null"`

	Location       string
	AdditionalInfo string

	Facilities []VerificationProductionCenterFacility `gorm:"foreignKey:ProductionCenterID"`

	CreatedAt time.Time
}

type VerificationProductionCenterFacility struct {
	ID                 uint `gorm:"primaryKey"`
	ProductionCenterID uint
	FacilityName       string

	CreatedAt time.Time
}