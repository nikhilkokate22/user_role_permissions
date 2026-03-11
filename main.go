package main

import (
	"log"
	app "user_role_permissions/app"

	"github.com/gin-gonic/gin"

	"user_role_permissions/config"
	"user_role_permissions/controller"
	"user_role_permissions/repository"
	"user_role_permissions/routes"
	"user_role_permissions/service"
)

func main() {

	cfg := config.LoadConfig()
	config.AppConfig = cfg
	config.InitLogger()

	db := config.ConnectDatabase(cfg)

	// repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	applicationRepo := repository.NewApplicantRepository(db)
	verificationRepo := repository.NewVerificationRepo(db)

	// services
	userService := service.NewUserService(userRepo, db, cfg)
	roleService := service.NewRoleService(roleRepo)
	menuService := service.NewMenuService(menuRepo, db)
	applicationService := service.NewApplicantService(applicationRepo, db)
	verificationService := service.NewVerificationService(verificationRepo, db)

	// controllers
	userController := controller.NewUserController(userService, db)
	roleController := controller.NewRolesController(roleService, db)
	permissionController := controller.NewPermissionController(permissionRepo, db)
	menuController := controller.NewMenuController(menuService)
	applicationController := controller.NewApplicantController(applicationService)
	verificationController := controller.NewVerificationController(verificationService)

	// create app container
	app := &app.App{
		DB: db,
		UserController: userController,
		RoleController: roleController,
		PermissionController: permissionController,
		PermissionRepo: permissionRepo,
		MenuController: menuController,
		ApplicationController: applicationController,
		VerificationController: verificationController,
	}

	r := gin.Default()

	routes.RegisterRoutes(r, app)

	log.Println("Server running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
