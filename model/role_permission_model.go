package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RolePermission struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      uuid.UUID      `gorm:"column:uuid;type:uuid;unique;not null"`
	RoleID    uint           `gorm:"column:role_id" json:"role_id"`
	MenuID    uint           `gorm:"column:menu_id" json:"menu_id"`
	CanCreate bool           `gorm:"column:can_create;default:false" json:"can_create"`
	CanRead   bool           `gorm:"column:can_read;default:false" json:"can_read"`
	CanUpdate bool           `gorm:"column:can_update;default:false" json:"can_update"`
	CanDelete bool           `gorm:"column:can_delete;default:false" json:"can_delete"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (rp RolePermission) TableName() string {
	return "role_permissions"
}