package controller

import (
	"user_role_permissions/dto"
	"user_role_permissions/service"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController interface {
	CreateRole(c *gin.Context)
	ListRoles(c *gin.Context)
	DeleteRole(c *gin.Context)
	UpdateRole(c *gin.Context) 
}

type roleController struct{
	roleService service.RoleService
	db *gorm.DB
}

func NewRolesController(roleService service.RoleService, db *gorm.DB) RoleController{
	return &roleController{
		roleService:roleService,
		db:db,
	}
}

func (rc *roleController) CreateRole(c *gin.Context) {

	var req dto.CreateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	currentRoleID := c.GetUint("role_id")

	err := rc.roleService.CreateRole(rc.db, currentRoleID, req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "Role created successfully")
}


func (rc *roleController) ListRoles(c *gin.Context) {

	var req dto.ListRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := rc.roleService.ListRoles(rc.db)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithData(c, "Role list fetched", resp)
}

func (rc *roleController) UpdateRole(c *gin.Context) {

	var req dto.UpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	currentRoleID := c.GetUint("role_id")

	err := rc.roleService.UpdateRole(rc.db, currentRoleID, req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "Role updated successfully")
}

func (rc *roleController) DeleteRole(c *gin.Context) {

	var req dto.DeleteRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	currentRoleID := c.GetUint("role_id")

	err := rc.roleService.DeleteRole(rc.db, currentRoleID, req.ID)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "Role deleted successfully")
}
