package repository

import (
	"user_role_permissions/model"

	"gorm.io/gorm"
)

type TurnoverRepository interface {
	DeleteByTraceID(tx *gorm.DB, traceID string) error
	CreateTurnover(tx *gorm.DB, turnover *model.TurnoverDetails) error
	CreateYear(tx *gorm.DB, year *model.TurnoverYear) error
	// CreateDocument(tx *gorm.DB, doc *model.TurnoverDocument) error
}

type turnoverRepository struct {
	DB *gorm.DB
}

func NewTurnoverRepository(db *gorm.DB) TurnoverRepository {
	return &turnoverRepository{DB: db}
}

func (r *turnoverRepository) DeleteByTraceID(tx *gorm.DB, traceID string) error {

	// First delete child records (if no cascade)
	if err := tx.Where("turnover_id IN (?)",
		tx.Model(&model.TurnoverDetails{}).
			Select("id").
			Where("trace_id = ?", traceID),
	).Delete(&model.TurnoverYear{}).Error; err != nil {
		return err
	}

	// Then delete parent
	if err := tx.Where("trace_id = ?", traceID).
		Delete(&model.TurnoverDetails{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *turnoverRepository) CreateTurnover(tx *gorm.DB, turnover *model.TurnoverDetails) error {
	return tx.Create(turnover).Error
}

func (r *turnoverRepository) CreateYear(tx *gorm.DB, year *model.TurnoverYear) error {
	return tx.Create(year).Error
}

// func (r *turnoverRepository) CreateDocument(tx *gorm.DB, doc *model.TurnoverDocument) error {
// 	return tx.Create(doc).Error
// }
