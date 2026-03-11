package model

import "time"

type VerificationMaster struct {
	ID               uint      `gorm:"primaryKey"`
	VerificationUUID string    `gorm:"type:char(36);unique;not null"`
	ApplicationUUID  string    `gorm:"type:char(36);not null"`
	ApplicationNo    int
	DateOfReceipt    *time.Time
	ProductCodeNo    string
	DateOfVisit      *time.Time
	TcOfficeID       int
	CurrentStep      string

	Members            []VerificationMember           `gorm:"foreignKey:VerificationUUID;references:VerificationUUID"`
	ApplicantDetails   VerificationApplicantDetails   `gorm:"foreignKey:VerificationUUID;references:VerificationUUID"`
	ProductionDetails  []VerificationProductionDetail `gorm:"foreignKey:VerificationUUID;references:VerificationUUID"`
	TurnoverModes      []VerificationTurnoverMode     `gorm:"foreignKey:VerificationUUID;references:VerificationUUID"`
	ProductDetails     []VerificationProductDetail    `gorm:"foreignKey:VerificationUUID;references:VerificationUUID"`
	ProductionCenters  []VerificationProductionCenter `gorm:"foreignKey:VerificationUUID;references:VerificationUUID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type VerificationMember struct {
	ID               uint   `gorm:"primaryKey"`
	VerificationUUID string `gorm:"type:char(36);not null"`
	MemberID         int

	CreatedAt time.Time
}

func (m VerificationMaster) TableName() string {
	return "verification_master"
}

func (m VerificationMember) TableName() string {
	return "verification_members"
}