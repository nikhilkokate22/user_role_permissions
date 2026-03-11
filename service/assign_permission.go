package service

import (
	"errors"
	"user_role_permissions/dto"
	"user_role_permissions/repository"

	"gorm.io/gorm"
)

type AssignPermissionService interface {
	AssignPermissions(db *gorm.DB, req dto.AssignMenuPermissionRequest) error
	CheckPermission(roleID uint, menuID uint, action string) (bool, error)
}

type assignPermissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewAssignPermissionService(permissionRepo repository.PermissionRepository) AssignPermissionService {
	return &assignPermissionService{
		permissionRepo: permissionRepo,
	}
}

func (s *assignPermissionService) AssignPermissions(db *gorm.DB, req dto.AssignMenuPermissionRequest) error {

	return db.Transaction(func(tx *gorm.DB) error {

		if _, err := s.permissionRepo.GetRoleByID(tx, uint(req.Id)); err != nil {
			return errors.New("role not found")
		}

		if len(req.Menus) == 0 {
			return errors.New("menus cannot be empty")
		}

		if err := s.permissionRepo.AssignRolePermissions(tx, uint(req.Id), req.Menus); err != nil {
			return err
		}

		s.permissionRepo.ClearRoleCache(uint(req.Id))


		return nil
	})
}


func (s *assignPermissionService) CheckPermission(roleID uint, menuID uint, action string) (bool, error) {

	rp, err := s.permissionRepo.FindByRoleAndMenu(roleID, menuID)
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