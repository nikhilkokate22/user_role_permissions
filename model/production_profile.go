package model

import "time"

type ProductionProfile struct {
	ID      uint   `gorm:"column:id;primaryKey"`
	TraceID string `gorm:"column:trace_id;type:char(36);index"`

	ProductID            uint    `gorm:"column:product_id"`
	ProductDescription   string  `gorm:"column:product_description"`
	ProductSpecification string  `gorm:"column:product_specification"`
	Woolen               float64 `gorm:"column:woolen"`
	Blended              float64 `gorm:"column:blended"`
	Others               float64 `gorm:"column:others"`
	Total                float64 `gorm:"column:total"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (p ProductionProfile) TableName() string {
	return "production_profile"
}