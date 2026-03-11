package app

import (
	"user_role_permissions/controller"
	"user_role_permissions/repository"

	"gorm.io/gorm"
)

type App struct {
	DB *gorm.DB
	PermissionRepo repository.PermissionRepository
	UserController controller.UserController
	RoleController controller.RoleController
	PermissionController controller.PermissionController
	MenuController controller.MenuController

	ApplicationController controller.ApplicationController
	VerificationController controller.VerificationController
}
