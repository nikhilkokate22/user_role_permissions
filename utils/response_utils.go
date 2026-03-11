package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
type MiddlewareResponseBody struct {
	Message    string
	StatusCode int
	Body       struct{}
}

type ResponseBody struct {
	Message    string `json:"message"`
	StatusCode int    `json:"Status_code"`
	// DevMessage error       `json:"-"`
	DevMessage string      `json:"dev_message,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

func sendJSONResponse(c *gin.Context, statusCode int, message string, body interface{}) {
	response := ResponseBody{
		Message:    message,
		StatusCode: statusCode,
		Body:       body,
	}
	c.JSON(statusCode, response)
}

func abortWithJSON(c *gin.Context, statusCode int, message string, body interface{}) {
	response := ResponseBody{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
	}
	c.AbortWithStatusJSON(statusCode, response)
}

func ValidationResponse(c *gin.Context, message string) {
	// sendJSONResponse(c, http.StatusUnprocessableEntity, message, map[string]interface{}{})
	sendJSONResponse(c, http.StatusUnprocessableEntity, message, map[string]interface{}{})
}

func Success(c *gin.Context, message string) {
	c.JSON(200, gin.H{
		"success": true,
		"message": message,
	})
}

func SuccessWithData(c *gin.Context, message string, data interface{}) {
	c.IndentedJSON(200, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.JSON(400, gin.H{
		"success": false,
		"message": message,
	})
}

func Unauthorized(c *gin.Context) {
	c.JSON(401, gin.H{
		"success": false,
		"message": "unauthorized",
	})
}

func Forbidden(c *gin.Context) {
	c.JSON(403, gin.H{
		"success": false,
		"message": "forbidden",
	})
}

func InternalServerError(c *gin.Context, message string) {
	c.JSON(500, gin.H{
		"success": false,
		"message": message,
	})
}

func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": message,
	})
}


func RespondWithError(c *gin.Context, code int, message string) {
	response := MiddlewareResponseBody{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
	c.AbortWithStatusJSON(code, response)
}

func BadRequestAbortWithJSON(c *gin.Context, message string) {
	abortWithJSON(c, http.StatusBadRequest, message, nil)
}

func ForbiddenResponse(c *gin.Context, message string) {
	sendJSONResponse(c, http.StatusForbidden, message, nil)
}

func SuccessResponse(c *gin.Context, message string, body interface{}) {
	sendJSONResponse(c, http.StatusOK, message, body)
}

func InternalServerErrorResponse(c *gin.Context, err error) {
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		DevMessage: err.Error(),
	}
	c.JSON(http.StatusInternalServerError, response)
}

func PermissionDenied(c *gin.Context, roleID uint, menuID uint) {
	// c.JSON(403, gin.H{
	// 	"status":  403,
	// 	"message": "Permission denied - action not allowed",
	// 	"role_id": roleID,
	// 	"action":  action,
	// 	"menu":    menu,
	// })
	response := ResponseBody{
		StatusCode: http.StatusForbidden,
		Message: "Permission denied - action not allowed",
	}
	c.JSON(http.StatusForbidden, response)
	c.Abort()
}