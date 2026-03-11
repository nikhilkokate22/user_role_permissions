package middleware

import (
	"fmt"
	"user_role_permissions/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)
func ValidateToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	logrus.Info("ValidateToken@ entered into validateToken func")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		
		return []byte(config.AppConfig.JwtSecret), nil
	})


	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, fmt.Errorf("invalid token")
	}
	logrus.Info("ValidateToken@ token validated")

	return token, claims, nil
}
