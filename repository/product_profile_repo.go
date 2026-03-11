package repository

import (
	"user_role_permissions/model"
	"gorm.io/gorm"
)

type ProductionProfileRepository interface {
	BulkInsert(items []model.ProductionProfile) error
}

type productionProfileRepository struct {
	db *gorm.DB
}

func NewProductionProfileRepository(db *gorm.DB) ProductionProfileRepository {
	return &productionProfileRepository{db}
}

func (r *productionProfileRepository) BulkInsert(items []model.ProductionProfile) error {
	return r.db.Create(&items).Error
}
