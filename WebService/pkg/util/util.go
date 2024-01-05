package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/config"
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
	user := c.Get("user")

	// 若尚未登入
	if user == nil {
		return nil
	}

	jwtToken := user.(*jwt.Token)
	claims := jwtToken.Claims.(*JwtCustomClaims)
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

func SendJSONBadRequest(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, echo.Map{
		"message": "請求參數錯誤！",
	})
}

func SendJSONInternalServerError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, echo.Map{
		"message": "系統發生錯誤！",
	})
}

func SendJSONResponse(
	c echo.Context,
	microserviceResponse protoreflect.ProtoMessage,
) error {
	errorMessage := "SendJSONByProtojson failed! error: %w"

	result, err := protojson.MarshalOptions{
		EmitUnpopulated: true, // Zero value 的欄位不要省略
	}.Marshal(microserviceResponse)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return SendJSONInternalServerError(c)
	}

	return c.JSONBlob(http.StatusOK, result)
}
