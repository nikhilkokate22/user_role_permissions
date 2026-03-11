package model

import "time"

type DeclarationDetails struct {
	ID                  uint       `gorm:"column:id;primaryKey"`
	TraceID             string     `gorm:"column:trace_id;type:varchar(100);uniqueIndex;not null"`
	UndertakingName     string     `gorm:"column:undertaking_name;type:varchar(255)"`
	UndertakingDate     *time.Time `gorm:"column:undertaking_date"`
	UndertakingPlace    string     `gorm:"column:undertaking_place;undertaking_place;type:varchar(255)"`
	UndertakingAccepted bool       `gorm:"column:undertaking_accepted"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

func (d DeclarationDetails) TableName() string {
	return "declaration"
}