package repository

import (
	"user_role_permissions/model"
	"gorm.io/gorm"
)

type DeclarationRepository interface{
	SaveDeclaration(data *model.DeclarationDetails) error
}

type declarationRepository struct {
	DB *gorm.DB
}

func NewDeclarationRepository(db *gorm.DB) DeclarationRepository {
	return &declarationRepository{DB: db}
}

func (r *declarationRepository) SaveDeclaration(data *model.DeclarationDetails) error {
	return r.DB.
		Where("trace_id = ?", data.TraceID).
		Save(data).Error
}
