package middleware

import (
	"strings"
	"user_role_permissions/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {

		roleAny, ok := c.Get("role_name")
		if !ok {
			logrus.Warn("SuperAdminOnly@ role_name missing in context")
			utils.Forbidden(c)
			c.Abort()
			return
		}

		role, ok := roleAny.(string)
		if !ok {
			logrus.Error("SuperAdminOnly@ role_name type assertion failed")
			utils.Forbidden(c)
			c.Abort()
			return
		}

		logrus.Infof("SuperAdminOnly@ role=%s attempting access", role)

		if strings.ToLower(role) == "super_admin" {
			logrus.Info("SuperAdminOnly@ access granted")
			c.Next()
			return
		}

		logrus.Warnf("SuperAdminOnly@ access denied for role=%s", role)
		utils.Forbidden(c)
		c.Abort()
	}
}
