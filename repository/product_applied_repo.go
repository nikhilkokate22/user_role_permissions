package repository

import (
	"user_role_permissions/model"

	"gorm.io/gorm"
)

type ProductAppliedRepository interface{
	SaveProducts(products []model.ProductApplied) error
	SavePayment(payment *model.ProductPaymentDetails) error
}

type productAppliedRepository struct {
	DB *gorm.DB
}

func NewProductAppliedRepository(db *gorm.DB) ProductAppliedRepository {
	return &productAppliedRepository{DB: db}
}

func (r *productAppliedRepository) SaveProducts(products []model.ProductApplied) error {
	return r.DB.Create(&products).Error
}

func (r *productAppliedRepository) SavePayment(payment *model.ProductPaymentDetails) error {
	return r.DB.
		Where("trace_id = ?", payment.TraceID).
		Save(payment).Error
}
