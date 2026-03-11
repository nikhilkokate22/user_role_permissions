package controller

import (
	"user_role_permissions/dto"
	"user_role_permissions/service"

	"github.com/gin-gonic/gin"
)

type ApplicationController interface {
	HandleStep(c *gin.Context) 
	ListApplications(c *gin.Context)
}

type applicationController struct{
	service service.ApplicationService

} 

func NewApplicantController(service service.ApplicationService) ApplicationController {
	return &applicationController{service}
}

func (ac *applicationController) HandleStep(c *gin.Context) {

	step := c.PostForm("current_step")
	traceID := c.PostForm("trace_id")

	newTraceID, err := ac.service.ProcessStep(c, step, traceID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if step == "applicant_details" {
		c.JSON(200, gin.H{
			"message":  "Step processed",
			"trace_id": newTraceID,
		})
		return
	}

	c.JSON(200, gin.H{"message": "Step processed"})
}

func (ac *applicationController) ListApplications(c *gin.Context) {

	var req dto.ApplicationListRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	data, total, filtered_count, err := ac.service.ListApplications(
		req.Limit,
		req.Offset,
		req.ApplicantName,
	)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Applications fetched successfully",
		"meta": gin.H{
			"total_count": total,
			"filter_count":  filtered_count,
		},
		"data": data,
	})
}