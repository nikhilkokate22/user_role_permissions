package middleware

import (
	"strings"
	"user_role_permissions/repository"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logrus.Warn("AuthMiddleware@ missing token")
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		_, claims, err := ValidateToken(tokenStr)
		if err != nil {
			logrus.Error("AuthMiddleware@ invalid token")
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		logrus.Info("AuthMiddleware@ Token validate successfuly")

		c.Set("user_id", uint(claims["uid"].(float64)))
		c.Set("role_id", uint(claims["role_id"].(float64)))

		if roleName, ok := claims["role_name"].(string); ok {
			c.Set("role_name", roleName)
		}

		c.Next()
	}
}

func Permission(repo repository.PermissionRepository, menuID uint, action string) gin.HandlerFunc {

	return func(c *gin.Context) {

		roleID := c.GetUint("role_id")

		if roleID == 0 {
			logrus.Warn("Permission@ Missing role_id in context")
			c.AbortWithStatus(401)
			return
		}

		logrus.Infof("Permission@ Checking role=%d menu=%d action=%s", roleID, menuID, action)

		allowed, err := repo.CheckAction(roleID, menuID, action)
		if err != nil {
			logrus.Warnf("Permission@ No permission record found role_id=%d menu_id=%d", roleID, menuID)
			c.AbortWithStatus(403)
			return
		}

		if !allowed {
			logrus.Warnf("Permission@ Action not allowed role_id=%d action=%s", roleID, action)
			utils.PermissionDenied(c, roleID, menuID)
			return
		}
		logrus.Infof("Permission@ Permission granted role_id=%d", roleID)

		c.Next()

	}
}
