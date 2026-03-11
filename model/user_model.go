package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"column:id;primaryKey"`
	UUID      string         `gorm:"column:uuid;type:uuid;unique;not null"`
	Name      string         `gorm:"column:name;type:varchar(100)" json:"name"`
	Email     string         `gorm:"column:email;type:text;not null"`
	Mobile    string         `gorm:"column:mobile;type:text"`
	Password  string         `gorm:"column:password;type:text" json:"-"`
	RoleID    uint           `gorm:"column:role_id" json:"role_id"`
	Role      Role           `gorm:"foreignKey:RoleID;references:ID"`
	CreatedBy uint           `gorm:"column:created_by" json:"created_by"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (u User) TableName() string {
	return "users"
}
