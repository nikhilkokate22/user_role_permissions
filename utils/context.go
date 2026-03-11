package utils

import "github.com/gin-gonic/gin"

func GetUserIDFromContext(c *gin.Context) uint {
	if val, exists := c.Get("user_id"); exists {
		return val.(uint)
	}
	return 0
}

func GetRoleIDFromContext(c *gin.Context) uint {
	if val, exists := c.Get("role_id"); exists {
		return val.(uint)
	}
	return 0
}
