package routes

import (
	"github.com/gin-gonic/gin"

	"user_role_permissions/app"
	"user_role_permissions/middleware"
)

func RegisterRoutes(r *gin.Engine, app *app.App) {

	public := r.Group("/api")
	{
		public.POST("/login", app.UserController.UserLogin)
		public.POST("/users/register", app.UserController.CreateUser)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		roles := protected.Group("/roles")
		{
			roles.POST("/create", middleware.Permission(app.PermissionRepo, 2, "create"), app.RoleController.CreateRole)
			roles.POST("/list", middleware.Permission(app.PermissionRepo, 2, "read"), app.RoleController.ListRoles)
			roles.POST("/update", middleware.Permission(app.PermissionRepo, 2, "update"), app.RoleController.UpdateRole)
			roles.POST("/delete", middleware.Permission(app.PermissionRepo, 2, "delete"), app.RoleController.DeleteRole)
			roles.PUT("/:role_id/permissions", middleware.Permission(app.PermissionRepo, 2, "update"), app.PermissionController.AssignRolePermissions)

			roles.GET("/:role_id/permissions", middleware.Permission(app.PermissionRepo, 2, "read"), app.PermissionController.GetRolePermissions)

		}

		users := protected.Group("/users")
		{
			users.POST("/list", middleware.Permission(app.PermissionRepo, 1, "read"), app.UserController.ListUsers)

			users.POST("/update", middleware.Permission(app.PermissionRepo, 1, "update"), app.UserController.UpdateUser)

			users.POST("/delete", middleware.Permission(app.PermissionRepo, 1, "delete"), app.UserController.DeleteUser)


		}

		menus := protected.Group("/menu")
		{
			menus.POST("/create", middleware.SuperAdminOnly(), app.MenuController.CreateMenu)
			menus.GET("/list", middleware.SuperAdminOnly(), app.MenuController.GetMenuList)
			menus.GET("/:uuid", middleware.SuperAdminOnly(), app.MenuController.GetMenuByUUID)
			menus.PUT("/:uuid", middleware.SuperAdminOnly(), app.MenuController.UpdateMenu)
			menus.DELETE("/:uuid", middleware.SuperAdminOnly(), app.MenuController.DeleteMenu)
		}

		
	}

	application := r.Group("/application")
	{
		application.POST("/", app.ApplicationController.HandleStep)
		application.POST("/list", app.ApplicationController.ListApplications)
		application.POST("/verification", app.VerificationController.ProcessVerification)
	}	



}
