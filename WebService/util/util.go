package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"github.com/kakurineuin/learn-english-microservices/web-service/config"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GetJWTClaims(c echo.Context) *JwtCustomClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)
	return claims
}

func GetJWTToken(userId, username, role string) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		userId,
		username,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	signedToken, err := token.SignedString([]byte(config.EnvJWTSecretKey()))
	if err != nil {
		return "", fmt.Errorf("getJWTToken failed! error: %w", err)
	}

	return signedToken, nil
}

func SendJSONError(c echo.Context, httpErrorStatus int) error {
	return c.JSON(httpErrorStatus, echo.Map{
		"message": "系統發生錯誤！",
	})
}
