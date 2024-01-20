package main

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/config"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/handler/exam"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/handler/user"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/handler/word"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/examservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/wordservice"
	mymiddleware "github.com/kakurineuin/learn-english-microservices/web-service/pkg/middleware"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

func main() {
	// 讀取環境變數
	loadEnv()

	ctx := context.Background()

	// 連線到資料庫
	databaseRepository := repository.NewMongoDBRepository(config.EnvDatabaseName())
	err := databaseRepository.ConnectDB(ctx, config.EnvMongoDBURI())
	if err != nil {
		log.Fatal(err)
	}

	// 程式結束時，結束資料庫連線
	defer func() {
		if err := databaseRepository.DisconnectDB(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// 微服務
	// ExamService
	examService := examservice.New(config.EnvExamServiceServerAddress())
	err = examService.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// 程式結束時，結束微服務連線
	defer func() {
		if err := examService.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	// WordService
	wordService := wordservice.New(config.EnvWordServiceServerAddress())
	err = wordService.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// 程式結束時，結束微服務連線
	defer func() {
		if err := wordService.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()

	// Handlers
	userHandler := user.NewHandler(databaseRepository)
	examHandler := exam.NewHandler(examService, databaseRepository)
	wordHandler := word.NewHandler(wordService)

	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	// Middleware
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "操作逾時",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Errorf("The operation has timed out, path: %s", c.Path())
		},
		Timeout: 30 * time.Second,
	}))

	if config.EnvEnableCSRF() {
		e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			CookieName:  "XSRF-TOKEN",
			TokenLookup: "header:X-XSRF-TOKEN",
		}))
	}

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(30)))
	e.Use(mymiddleware.UserHistory(databaseRepository))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/assets", "frontend/dist/assets")
	e.File("/*", "frontend/dist/index.html")

	setupAPIHandlers(e, userHandler, examHandler, wordHandler)

	e.Logger.Info("Echo starts to listin at " + config.EnvWebServiceServerAddress())
	e.Logger.Fatal(e.Start(config.EnvWebServiceServerAddress()))
}

func setupAPIHandlers(
	e *echo.Echo,
	userHandler user.UserHandler,
	examHandler exam.ExamHandler,
	wordHandler word.WordHandler,
) {
	// API
	api := e.Group("/api")

	// 登入
	api.POST("/login", userHandler.Login)

	// 註冊
	api.POST("/user", userHandler.CreateUser)

	// 未登入時的 ExamInfo
	api.GET("/exam/info", examHandler.FindExamInfosWhenNotSignIn)

	// Restricted group，需要登入後才能呼叫的 API
	restrictedApi := api.Group("/restricted")

	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(util.JwtCustomClaims)
		},
		SigningKey: []byte(config.EnvJWTSecretKey()),
	}
	restrictedApi.Use(echojwt.WithConfig(config))

	// ExamInfo
	restrictedApi.GET("/exam/info", examHandler.FindExamInfosWhenSignIn)

	// Exam
	restrictedApi.GET("/exam", examHandler.FindExams)
	restrictedApi.POST("/exam", examHandler.CreateExam)
	restrictedApi.PATCH("/exam", examHandler.UpdateExam)
	restrictedApi.DELETE("/exam/:examId", examHandler.DeleteExam)
	restrictedApi.GET("/exam/:examId/question", examHandler.FindQuestions)
	restrictedApi.POST("/exam/:examId/question", examHandler.CreateQuestion)
	restrictedApi.PATCH("/exam/:examId/question", examHandler.UpdateQuestion)
	restrictedApi.DELETE("/exam/:examId/question/:questionId", examHandler.DeleteQuestion)
	restrictedApi.GET("/exam/:examId/start", examHandler.FindRandomQuestions)
	restrictedApi.POST("/exam/:examId/record", examHandler.CreateExamRecord)
	restrictedApi.GET("/exam/:examId/record/overview", examHandler.FindExamRecordOverview)
	restrictedApi.GET("/exam/:examId/record", examHandler.FindExamRecords)

	// Word
	restrictedApi.GET("/word/:word", wordHandler.FindWordMeanings)
	restrictedApi.POST("/word/favorite", wordHandler.CreateFavoriteWordMeaning)
	restrictedApi.GET("/word/favorite", wordHandler.FindFavoriteWordMeanings)
	restrictedApi.DELETE(
		"/word/favorite/:favoriteWordMeaningId",
		wordHandler.DeleteFavoriteWordMeaning,
	)
	restrictedApi.GET("/word/card", wordHandler.FindRandomFavoriteWordMeanings)

	// User
	restrictedApi.GET("/user/history", userHandler.FindUserHistories)
}

func loadEnv() {
	env := os.Getenv("WEB_SERVICE_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}
