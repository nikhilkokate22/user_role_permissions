package controller

import (
	"errors"
	"net/http"
	"user_role_permissions/dto"
	"user_role_permissions/service"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserController interface {
	CreateUser(c *gin.Context)
	ListUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	UserLogin(c *gin.Context)
}

type userController struct {
	userService service.UserService
	db          *gorm.DB
}

func NewUserController(userService service.UserService, db *gorm.DB) UserController {
	return &userController{
		userService: userService,
		db:          db,
	}
}

func (uc *userController) CreateUser(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("CreateUser@panic:", panicInfo)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userID := utils.GetUserIDFromContext(c)

	if err := uc.userService.CreateUser(req, userID); err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "User created")
}

func (uc *userController) ListUsers(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UpdateUser panic:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var req dto.UserListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	resp, err := uc.userService.ListUsers(uc.db, req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.SuccessWithData(c, "User list fetched", resp)
}

func (u *userController) UpdateUser(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UpdateUser panic:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	err := u.userService.UpdateUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

func (uc *userController) DeleteUser(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("DeleteUser panic:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var req dto.DeleteUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if req.ID == 0 {
		utils.BadRequest(c, "invalid user id")
		return
	}

	err := uc.userService.DeleteUser(uc.db, req)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFoundResponse(c, "User not found")
			return
		}

		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Success(c, "User deleted successfully")
}


func (s *userController) UserLogin(c *gin.Context) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			logrus.Error("UserLogin panic:", panicInfo)
			utils.InternalServerErrorResponse(c, panicInfo.(error))
		}
	}()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Warn("UserLogin@ Invalid login request payload")
		utils.ValidationResponse(c, "invalid request payload")
		return
	}

	if validationResp := utils.ValidateRequest(c, req); validationResp != nil {
		logrus.Warn("Login request validation failed")
		utils.ValidationResponse(c, validationResp.(string))
		return
	}

	userResp, token, err := s.userService.Login(c, req)
	if err != nil {
		logrus.Warn("Login failed:", err)

		switch err.Error() {
		case "invalid credentials":
			utils.BadRequestAbortWithJSON(c, err.Error())
		case "access denied":
			utils.ForbiddenResponse(c, err.Error())
		default:
			utils.InternalServerErrorResponse(c, err)
		}
		return
	}
	logrus.Info("Login successful")
	utils.SuccessResponse(c, "Login successful", gin.H{
		"user":  userResp,
		"token": token,
	})
}
