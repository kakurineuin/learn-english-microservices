package main

import (
	"context"
	"net"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pb"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/config"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/endpoint"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/service"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/transport"
)

func main() {
	// 讀取環境變數
	loadEnv()

	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	errorLogger := level.Error(logger)

	ctx := context.Background()

	// 連線到資料庫
	databaseRepository := repository.NewMongoDBRepository(config.EnvDatabaseName())
	err := databaseRepository.ConnectDB(ctx, config.EnvMongoDBURI())
	if err != nil {
		errorLogger.Log("msg", "Connect DB fail", "err", err)
		os.Exit(1)
	}

	// 程式結束時，結束資料庫連線
	defer func() {
		if err := databaseRepository.DisconnectDB(ctx); err != nil {
			errorLogger.Log("msg", "Disconnect DB fail", "err", err)
			os.Exit(1)
		}
	}()

	listener, err := net.Listen("tcp", config.EnvExamServiceServerAddress())
	if err != nil {
		errorLogger.Log("msg", "net listen fail", "err", err)
		os.Exit(1)
	}

	examService := service.New(logger, databaseRepository)
	examEndpoints := endpoint.MakeEndpoints(examService, logger)
	myGrpcServer := transport.NewGRPCServer(examEndpoints, logger)

	grpcServer := grpc.NewServer()
	pb.RegisterExamServiceServer(grpcServer, myGrpcServer)
	reflection.Register(grpcServer)
	level.Info(logger).Log("msg", "Starting gRPC server at "+config.EnvExamServiceServerAddress())
	err = grpcServer.Serve(listener)

	if err != nil {
		errorLogger.Log("msg", "grpcServer serve fail", "err", err)
	}
}

func loadEnv() {
	env := os.Getenv("SERVICE_ENV")
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
