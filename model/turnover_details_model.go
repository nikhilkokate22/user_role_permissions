package model

import (
	"time"
)

type TurnoverDetails struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	TraceID   string         `gorm:"column:trace_id;type:char(36);not null;index"`
	ModeID    uint           `gorm:"column:mode_id;not null"`
	Document  string         `gorm:"column:document;type:varchar(500)"`
	Years     []TurnoverYear `gorm:"foreignKey:TurnoverID"`
	CreatedAt time.Time      `gorm:"column:created_at"`
}

type TurnoverYear struct {
	ID            uint      `gorm:"column:id;primaryKey"`
	TurnoverID    uint      `gorm:"column:turnover_id;not null;index"`
	FinancialYear string    `gorm:"column:financial_year;type:varchar(20);not null"`
	Amount        float64   `gorm:"column:amount;type:decimal(15,2)"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

// type TurnoverDocument struct {
// 	ID        uint      `gorm:"column:id;primaryKey"`
// 	TraceID   string    `gorm:"column:trace_id;index"`
// 	FilePath  string    `gorm:"column:file_path"`
// 	CreatedAt time.Time `gorm:"column:created_at"`

// }

// type ApplicationDetails struct {
// 	gorm.Model
// 	TraceID     string `gorm:"uniqueIndex"`
// 	CurrentStep string
// }
