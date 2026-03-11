package model

import "time"

type ProductionCenter struct {
	ID             uint                       `gorm:"primaryKey"`
	TraceID        string                     `gorm:"type:char(36);not null;index"`
	Location       string                     `gorm:"type:varchar(255)"`
	AdditionalInfo string                     `gorm:"type:text"`
	Facilities     []ProductionCenterFacility `gorm:"foreignKey:ProductionCenterID"`
	CreatedAt      time.Time                  `gorm:"column:created_at"`
	UpdatedAt      time.Time                  `gorm:"column:updated_at"`
}

type ProductionCenterFacility struct {
	ID                 uint      `gorm:"primaryKey"`
	ProductionCenterID uint      `gorm:"not null;index"`
	FacilityName       string    `gorm:"type:varchar(100);not null"`
	CreatedAt          time.Time `gorm:"column:created_at"`
}
	
func (p ProductionCenter) TableName() string {
	return "production_center"
}

func (p ProductionCenterFacility) TableName() string {
	return "production_center_facilities"
}
