package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/kakurineuin/learn-english-microservices/web-service/config"
	"github.com/kakurineuin/learn-english-microservices/web-service/database"
	"github.com/kakurineuin/learn-english-microservices/web-service/model/user"
	"github.com/kakurineuin/learn-english-microservices/web-service/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type JwtCustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func CreateUser(c echo.Context) error {
	jsonMap, err := util.GetJSONBody(c)
	if err != nil {
		return fmt.Errorf("CreateUser failed! %w", err)
	}

	username := jsonMap["username"].(string)

	// 檢查使用者名稱是否已被註冊
	var findUser user.User
	usersCollection := database.GetCollection("users")
	err = usersCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&findUser)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "此使用者名稱已被註冊",
		})
	} else if err != mongo.ErrNoDocuments {
		return fmt.Errorf("CreateUser check user failed! %w", err)
	}

	password := jsonMap["password"].(string)

	// Password Hashing
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return fmt.Errorf("CreateUser hash password failed! %w", err)
	}

	encryptedPassword := string(bytes)
	user := &user.User{
		Username: username,
		Password: encryptedPassword,
		Role:     "user",
	}
	_, dbErr := usersCollection.InsertOne(context.TODO(), user)
	if dbErr != nil {
		log.Error(err)
		return fmt.Errorf("CreateUser failed! error: %w", dbErr)
	}

	// 產生 JWT
	token, err := getJWTToken(username, user.Role)
	if err != nil {
		return fmt.Errorf("CreateUser get JWT failed! error: %w", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func Login(c echo.Context) error {
	jsonMap, err := util.GetJSONBody(c)
	if err != nil {
		return fmt.Errorf("Login failed! %w", err)
	}

	username := jsonMap["username"].(string)
	password := jsonMap["password"].(string)

	// Check user by MongoDB
	var user user.User
	usersCollection := database.GetCollection("users")
	err = usersCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&user)
	if err != nil {

		// 查無此帳號，表示使用者輸入錯誤的帳號
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "帳號錯誤",
			})
		}

		// DB error
		return fmt.Errorf("login failed! error: %w", err)
	}

	// 檢查密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "密碼錯誤",
		})
	}

	// 產生 JWT
	token, err := getJWTToken(username, user.Role)
	if err != nil {
		return fmt.Errorf("login failed! error: %w", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func getJWTToken(username, role string) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
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
