package main

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/kakurineuin/learn-english-microservices/web-service/config"
	"github.com/kakurineuin/learn-english-microservices/web-service/database"
	"github.com/kakurineuin/learn-english-microservices/web-service/handler/user"
	"github.com/kakurineuin/learn-english-microservices/web-service/handler/word"
	"github.com/kakurineuin/learn-english-microservices/web-service/microservice"
)

// func restricted(c echo.Context) error {
// 	token := c.Get("user").(*jwt.Token)
// 	claims := token.Claims.(*user.JwtCustomClaims)
// 	username := claims.Username
// 	role := claims.Role
// 	return c.String(http.StatusOK, "Welcome "+username+"! role: "+role)
// }

func main() {
	loadEnv()

	// 連線到資料庫
	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	// 程式結束時，結束資料庫連線
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			log.Fatal(err)
		}
	}()

	// 連線到微服務
	if err := microservice.Connect(); err != nil {
		log.Fatal(err)
	}

	// 程式結束時，結束微服務連線
	defer func() {
		if err := microservice.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets", "frontend/dist/assets")
	e.File("/*", "frontend/dist/index.html")

	setupAPIHandlers(e)

	e.Logger.Info("Echo starts to listin...")
	e.Logger.Fatal(e.Start(":1323"))
}

func loadEnv() {
	switch os.Getenv("LEARN_ENGLISH_ENV") {
	case "PROD":
		godotenv.Load(".env.production")
	default:
		godotenv.Load(".env.local")
	}
}

func home(c echo.Context) error {
	return c.String(http.StatusOK, "")
}

func setupAPIHandlers(e *echo.Echo) {
	// API
	api := e.Group("/api")

	// 登入
	api.POST("/login", user.Login)

	// 註冊
	api.POST("/user", user.CreateUser)

	// Restricted group，需要登入後才能呼叫的 API
	restrictedApi := api.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(user.JwtCustomClaims)
		},
		SigningKey: []byte(config.EnvJWTSecretKey()),
	}
	restrictedApi.Use(echojwt.WithConfig(config))

	restrictedApi.GET("/word/:word", word.FindWordMeanings)
}
