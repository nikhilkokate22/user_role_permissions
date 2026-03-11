package repository

import (
	"gorm.io/gorm"
	"user_role_permissions/model"
)

type ProductionCenterRepository interface{
	Create(tx *gorm.DB, center *model.ProductionCenter) error
}

type productionCenterRepository struct {
	DB *gorm.DB
}

func NewProductionCenterRepository(db *gorm.DB) ProductionCenterRepository {
	return &productionCenterRepository{DB: db}
}

func (r *productionCenterRepository) DeleteByTraceID(tx *gorm.DB, traceID string) error {
	return tx.Where("trace_id = ?", traceID).
		Delete(&model.ProductionCenter{}).Error
}

func (r *productionCenterRepository) Create(tx *gorm.DB, center *model.ProductionCenter) error {
	return tx.Create(center).Error
}
