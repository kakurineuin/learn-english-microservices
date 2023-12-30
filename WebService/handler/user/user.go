package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/kakurineuin/learn-english-microservices/web-service/database"
	"github.com/kakurineuin/learn-english-microservices/web-service/model/user"
	"github.com/kakurineuin/learn-english-microservices/web-service/util"
)

func CreateUser(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	errorMessage := "CreateUser failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "請求參數不正確",
		})
	}

	username := requestBody.Username
	password := requestBody.Password

	// 檢查使用者名稱是否已被註冊
	var findUser user.User
	usersCollection := database.GetCollection("users")
	err := usersCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&findUser)
	if err == nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "此使用者名稱已被註冊",
		})
	} else if err != mongo.ErrNoDocuments {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "系統發生錯誤",
		})
	}

	// Password Hashing
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "系統發生錯誤",
		})
	}

	encryptedPassword := string(bytes)
	user := &user.User{
		Username: username,
		Password: encryptedPassword,
		Role:     "user",
	}
	_, err = usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "系統發生錯誤",
		})
	}

	// 產生 JWT
	token, err := util.GetJWTToken(user.Id.Hex(), username, user.Role)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "系統發生錯誤",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func Login(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	errorMessage := "Login failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "請求參數不正確",
		})
	}

	username := requestBody.Username
	password := requestBody.Password

	// Check user by MongoDB
	var user user.User
	usersCollection := database.GetCollection("users")
	err := usersCollection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&user)
	if err != nil {

		// 查無此帳號，表示使用者輸入錯誤的帳號
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "帳號錯誤",
			})
		}

		// DB error
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "系統發生錯誤",
		})
	}

	// 檢查密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "密碼錯誤",
		})
	}

	// 產生 JWT
	token, err := util.GetJWTToken(user.Id.Hex(), username, user.Role)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "系統發生錯誤",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
