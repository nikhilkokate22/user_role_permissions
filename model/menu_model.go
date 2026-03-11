package model

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      string         `gorm:"column:uuid"`
	Name      string         `gorm:"column:name"`
	Path      string         `gorm:"column:path"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (m Menu) TableName() string {
	return "menus"
}
