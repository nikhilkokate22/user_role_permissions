package repository

import (
	"user_role_permissions/dto"
	"user_role_permissions/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetPermission(roleID uint, menuID uint) (*model.RolePermission, error)
	GetRoleByID(db *gorm.DB, id uint) (model.Role, error)
	AssignRolePermissions(tx *gorm.DB, roleID uint, menus []dto.AssignMenu) error
	ClearRoleCache(roleID uint)

	GetPermissionsByRole(roleID uint) ([]model.RolePermission, error)
	UpdateRolePermission(tx *gorm.DB, rp *model.RolePermission) error
	CheckAction(roleID uint, menuID uint, action string) (bool, error)
	FindByRoleAndMenu(roleID uint, menuID uint) (*model.RolePermission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

var permissionCache = make(map[uint]map[uint]model.RolePermission)

func (r *permissionRepository) GetPermission(roleID uint, menuID uint) (*model.RolePermission, error) {

	if rolePerms, ok := permissionCache[roleID]; ok {
		if p, ok := rolePerms[menuID]; ok {
			cp := p
			return &cp, nil
		}

	}

	var rp model.RolePermission

	err := r.db.Where("role_id = ? AND menu_id = ?", roleID, menuID).First(&rp).Error
	// Joins("JOIN roles ON users.role_id = roles.id AND roles.deleted_at IS NULL")

	if err != nil {
		return nil, err
	}

	if _, ok := permissionCache[roleID]; !ok {
		permissionCache[roleID] = make(map[uint]model.RolePermission)
	}

	permissionCache[roleID][menuID] = rp

	return &rp, nil
}

func (r *permissionRepository) GetRoleByID(db *gorm.DB, id uint) (model.Role, error) {
	var role model.Role
	err := db.First(&role, id).Error // GORM auto filters deleted_at
	return role, err
}

func (r *permissionRepository) AssignRolePermissions(tx *gorm.DB, roleID uint, menus []dto.AssignMenu) error {

	if err := tx.
		Where("role_id = ?", roleID).
		Delete(&model.RolePermission{}).Error; err != nil {
		return err
	}

	var inserts []model.RolePermission

	for _, m := range menus {

		rp := model.RolePermission{
			UUID:   uuid.New(),
			RoleID: roleID,
			MenuID: uint(m.Id),
		}

		// Convert permission IDs to bool flags
		for _, p := range m.Permissions {

			switch p {
			case 1:
				rp.CanCreate = true
			case 2:
				rp.CanRead = true
			case 3:
				rp.CanUpdate = true
			case 4:
				rp.CanDelete = true
			}
		}

		inserts = append(inserts, rp)
	}

	if len(inserts) > 0 {
		if err := tx.Create(&inserts).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *permissionRepository) ClearRoleCache(roleID uint) {
	delete(permissionCache, roleID)
}

func (r *permissionRepository) GetPermissionsByRole(roleID uint) ([]model.RolePermission, error) {

	var perms []model.RolePermission

	err := r.db.
		Where("role_id = ?", roleID).
		Find(&perms).Error

	if err != nil {
		return nil, err
	}

	// Refresh cache completely
	roleMap := make(map[uint]model.RolePermission)
	for _, p := range perms {
		roleMap[p.MenuID] = p
	}
	permissionCache[roleID] = roleMap

	return perms, nil
}

func (r *permissionRepository) UpdateRolePermission(tx *gorm.DB, rp *model.RolePermission) error {

	if err := tx.Save(rp).Error; err != nil {
		return err
	}

	// Clear cache for safety
	r.ClearRoleCache(rp.RoleID)

	return nil
}

func (r *permissionRepository) CheckAction(roleID uint, menuID uint, action string) (bool, error) {

	rp, err := r.GetPermission(roleID, menuID)
	if err != nil {
		return false, err
	}

	switch action {
	case "create":
		return rp.CanCreate, nil
	case "read":
		return rp.CanRead, nil
	case "update":
		return rp.CanUpdate, nil
	case "delete":
		return rp.CanDelete, nil
	default:
		return false, nil
	}
}

func (r *permissionRepository) FindByRoleAndMenu(roleID uint, menuID uint) (*model.RolePermission, error) {
	var rp model.RolePermission
	err := r.db.Where("role_id = ? AND menu_id = ?", roleID, menuID).First(&rp).Error
	return &rp, err
}