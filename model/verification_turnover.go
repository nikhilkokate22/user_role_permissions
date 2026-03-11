package model

import "time"

type VerificationTurnoverMode struct {
	ID               uint                       `gorm:"column:id;primaryKey"`
	VerificationUUID string                     `gorm:"column:verification_uuid;type:char(36);not null;index"`
	ModeID           uint                       `gorm:"column:mode_id;not null"`
	Document         string                     `gorm:"column:document;type:varchar(500)"`
	Years            []VerificationTurnoverYear `gorm:"foreignKey:TurnoverModeID"`
	CreatedAt        time.Time
}

type VerificationTurnoverYear struct {
	ID             uint `gorm:"primaryKey"`
	TurnoverModeID uint
	Year           string
	Amount         float64

	CreatedAt time.Time
}
