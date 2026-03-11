package controller

import (
	"net/http"
	"user_role_permissions/dto"
	"user_role_permissions/service"

	"github.com/gin-gonic/gin"
)

type MenuController interface{
	CreateMenu(ctx *gin.Context) 
	GetMenuList(ctx *gin.Context)
	GetMenuByUUID(ctx *gin.Context)
	UpdateMenu(ctx *gin.Context)
	DeleteMenu(ctx *gin.Context)
}

type menuController struct {
	service service.MenuService
}

func NewMenuController(service service.MenuService) MenuController {
	return &menuController{
		service: service,
	}
}

func (c *menuController) CreateMenu(ctx *gin.Context) {

	var req dto.CreateMenuRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.service.CreateMenu(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "menu created successfully",
	})
}

func (c *menuController) GetMenuList(ctx *gin.Context) {

	var req dto.MenuListRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := c.service.GetMenuList(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *menuController) GetMenuByUUID(ctx *gin.Context) {

	uuid := ctx.Param("uuid")

	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "uuid is required",
		})
		return
	}

	menu, err := c.service.GetMenuByUUID(uuid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "menu not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, menu)
}

func (c *menuController) UpdateMenu(ctx *gin.Context) {

	uuid := ctx.Param("uuid")

	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "uuid is required",
		})
		return
	}

	var req dto.UpdateMenuRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := c.service.UpdateMenu(uuid, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "menu updated successfully",
	})
}

func (c *menuController) DeleteMenu(ctx *gin.Context) {

	uuid := ctx.Param("uuid")

	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "uuid is required",
		})
		return
	}

	if err := c.service.DeleteMenu(uuid); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "menu deleted successfully",
	})
}
