package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

var utilGetJWTToken = util.GetJWTToken

type UserHandler interface {
	CreateUser(c echo.Context) error
	Login(c echo.Context) error
}

type userHandler struct {
	databaseRepository repository.DatabaseRepository
}

func NewHandler(databaseRepository repository.DatabaseRepository) UserHandler {
	return &userHandler{
		databaseRepository: databaseRepository,
	}
}

func (handler userHandler) CreateUser(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	errorMessage := "CreateUser failed! error: %w"
	databaseRepository := handler.databaseRepository

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	username := requestBody.Username
	password := requestBody.Password

	// 檢查使用者名稱是否已被註冊
	findUser, err := databaseRepository.GetUserByUsername(context.TODO(), username)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}
	if findUser != nil {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "此使用者名稱已被註冊",
		})
	}

	// Password Hashing
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	encryptedPassword := string(bytes)
	user := model.User{
		Username: username,
		Password: encryptedPassword,
		Role:     "user",
	}
	userId, err := databaseRepository.CreateUser(context.TODO(), user)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	// 產生 JWT
	token, err := utilGetJWTToken(userId, username, user.Role)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (handler userHandler) Login(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	errorMessage := "Login failed! error: %w"
	databaseRepository := handler.databaseRepository

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	username := requestBody.Username
	password := requestBody.Password

	// Check user by MongoDB
	user, err := databaseRepository.GetUserByUsername(context.TODO(), username)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	fmt.Printf("================ user: %v", user)
	fmt.Println()

	// 查無此帳號，表示使用者輸入錯誤的帳號
	if user == nil {
		fmt.Println("================ 1")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "帳號錯誤",
		})
	}

	fmt.Println("================ 2")

	// 檢查密碼
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("================ 3")
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "密碼錯誤",
		})
	}

	// 產生 JWT
	token, err := utilGetJWTToken(user.Id.Hex(), username, user.Role)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
