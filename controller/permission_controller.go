package controller

import (
	"net/http"
	"strconv"

	"user_role_permissions/dto"
	"user_role_permissions/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PermissionController interface {
	AssignRolePermissions(c *gin.Context)
	GetRolePermissions(c *gin.Context)
}

type permissionController struct {
	repo repository.PermissionRepository
	db   *gorm.DB
}

func NewPermissionController(repo repository.PermissionRepository, db *gorm.DB) PermissionController {
	return &permissionController{
		repo: repo,
		db:   db,
	}
}

func (pc *permissionController) AssignRolePermissions(c *gin.Context) {

	roleIDParam := c.Param("role_id")
	roleIDInt, err := strconv.Atoi(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	roleID := uint(roleIDInt)

	var request []dto.AssignMenu

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate role exists
	role, err := pc.repo.GetRoleByID(pc.db, roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	logrus.Infof("AssignRolePermissions@ Role=%s ID=%d", role.Name, roleID)

	tx := pc.db.Begin()

	if err := pc.repo.AssignRolePermissions(tx, roleID, request); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}

	pc.repo.ClearRoleCache(roleID)

	c.JSON(http.StatusOK, gin.H{
		"message": "role permissions updated successfully",
	})
}

func (pc *permissionController) GetRolePermissions(c *gin.Context) {

	roleIDParam := c.Param("role_id")
	roleIDInt, err := strconv.Atoi(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role_id"})
		return
	}

	roleID := uint(roleIDInt)

	perms, err := pc.repo.GetPermissionsByRole(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, perms)
}
