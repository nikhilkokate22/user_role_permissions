package service

import (
	"errors"
	"user_role_permissions/dto"
	"user_role_permissions/model"
	"user_role_permissions/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleService interface {
	CreateRole(db *gorm.DB, currentRoleID uint, req dto.CreateRoleRequest) error
	UpdateRole(db *gorm.DB, currentRoleID uint, req dto.UpdateRoleRequest) error
	DeleteRole(db *gorm.DB, currentRoleID uint, roleID uint) error
	ListRoles(db *gorm.DB) ([]model.Role, error)
}

type roleService struct {
	rolesRepo repository.RolesRepository
}

func NewRoleService(rolesRepo repository.RolesRepository) RoleService {
	return &roleService{
		rolesRepo: rolesRepo,
	}
}

func (s *roleService) CreateRole(db *gorm.DB, currentRoleID uint, req dto.CreateRoleRequest) error {

	if err := s.validateSuperAdmin(db, currentRoleID); err != nil {
		return err
	}

	// check duplicate
	exists, _ := s.rolesRepo.GetRoleByName(db, req.Name)
	if exists.ID != 0 {
		return errors.New("role already exists")
	}

	role := model.Role{
		UUID: uuid.New(),
		Name: req.Name,
	}

	return s.rolesRepo.CreateRole(db, &role)
}

func (s *roleService) ListRoles(db *gorm.DB) ([]model.Role, error) {
	return s.rolesRepo.GetAllRoles(db)
}

func (s *roleService) UpdateRole(db *gorm.DB, currentRoleID uint, req dto.UpdateRoleRequest) error {

	if err := s.validateSuperAdmin(db, currentRoleID); err != nil {
		return err
	}

	role, err := s.rolesRepo.GetRoleByID(db, req.ID)
	if err != nil {
		return err
	}

	if role.Name == "SUPER_ADMIN" {
		return errors.New("super admin role cannot be modified")
	}

	role.Name = req.Name

	return s.rolesRepo.UpdateRole(db, &role)
}

func (s *roleService) validateSuperAdmin(db *gorm.DB, currentRoleID uint) error {

	role, err := s.rolesRepo.GetRoleByID(db, currentRoleID)
	if err != nil {
		return err
	}

	if role.Name != "SUPER_ADMIN" {
		return errors.New("only super admin can perform this action")
	}

	return nil
}

func (s *roleService) DeleteRole(db *gorm.DB, currentRoleID uint, roleID uint) error {

	return db.Transaction(func(tx *gorm.DB) error {

		if err := s.validateSuperAdmin(tx, currentRoleID); err != nil {
			return err
		}

		role, err := s.rolesRepo.GetRoleByID(tx, roleID)
		if err != nil {
			return errors.New("role not found")
		}

		if role.Name == "SUPER_ADMIN" {
			return errors.New("super admin role cannot be deleted")
		}

		count, err := s.rolesRepo.CountUsersByRoleID(tx, roleID)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.New("cannot delete role assigned to active users")
		}

		if err := s.rolesRepo.DeletePermissionsByRoleID(tx, roleID); err != nil {
			return err
		}

		if err := s.rolesRepo.DeleteRole(tx, roleID); err != nil {
			return err
		}

		return nil
	})
}

