package controller

import (
	"user_role_permissions/service"

	"github.com/gin-gonic/gin"
)

type VerificationController interface{
	ProcessVerification(c *gin.Context)
}

type verificationController struct{
	service service.VerificationService
}

func NewVerificationController(service service.VerificationService) VerificationController{
	return &verificationController{service: service}
}

func (s *verificationController) ProcessVerification(c *gin.Context) {

	step := c.PostForm("current_step")
	verificationID := c.PostForm("verification_id")

	resp, err := s.service.ProcessStep(c, step, verificationID)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}