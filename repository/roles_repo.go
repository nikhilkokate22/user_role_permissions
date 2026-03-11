package repository

import (
	"user_role_permissions/model"

	"gorm.io/gorm"
)

type RolesRepository interface {
	CreateRole(db *gorm.DB, role *model.Role) error
	GetRoleByID(db *gorm.DB, id uint) (model.Role, error)
	GetRoleByName(db *gorm.DB, name string) (model.Role, error)
	GetAllRoles(db *gorm.DB) ([]model.Role, error)
	UpdateRole(db *gorm.DB, role *model.Role) error
	DeleteRole(db *gorm.DB, id uint) error
	CountUsersByRoleID(db *gorm.DB, roleID uint) (int64, error)
	DeletePermissionsByRoleID(db *gorm.DB, roleID uint) error
	// AssignRolePermissions(tx *gorm.DB, roleID uint,	menus []dto.AssignMenu) error
}

type rolesRepository struct{}

func NewRoleRepository(db *gorm.DB) *rolesRepository {
	return &rolesRepository{}
}

func (r *rolesRepository) CreateRole(db *gorm.DB, role *model.Role) error {
	return db.Create(role).Error
}

func (r *rolesRepository) GetRoleByID(db *gorm.DB, id uint) (model.Role, error) {
	var role model.Role
	err := db.First(&role, id).Error // GORM auto filters deleted_at
	return role, err
}

func (r *rolesRepository) GetRoleByName(db *gorm.DB, name string) (model.Role, error) {
	var role model.Role
	err := db.Where("name = ?", name).First(&role).Error
	return role, err
}

func (r *rolesRepository) GetAllRoles(db *gorm.DB) ([]model.Role, error) {
	var roles []model.Role
	err := db.Order("id asc").Find(&roles).Error // auto filters deleted_at
	return roles, err
}

func (r *rolesRepository) UpdateRole(db *gorm.DB, role *model.Role) error {
	return db.Save(role).Error
}

func (r *rolesRepository) DeleteRole(db *gorm.DB, id uint) error {
	return db.Delete(&model.Role{}, id).Error // soft delete
}

func (r *rolesRepository) CountUsersByRoleID(db *gorm.DB, roleID uint) (int64, error) {
	var count int64
	err := db.Model(&model.User{}).
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Count(&count).Error
	return count, err
}

func (r *rolesRepository) DeletePermissionsByRoleID(db *gorm.DB, roleID uint) error {
	return db.
		Where("role_id = ?", roleID).
		Delete(&model.RolePermission{}).Error // soft delete if using gorm.DeletedAt
}

