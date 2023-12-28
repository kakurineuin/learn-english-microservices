package main

import (
	"net"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/kakurineuin/learn-english-microservices/word-service/pb"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/config"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/endpoint"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/service"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/transport"
)

const PORT = ":8090"

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	errorLogger := level.Error(logger)

	loadEnv()

	// 連線到資料庫
	databaseRepository := repository.NewMongoDBRepository("learnEnglish")
	err := databaseRepository.ConnectDB(config.EnvMongoDBURI())
	if err != nil {
		errorLogger.Log("msg", "Connect DB fail", "err", err)
		os.Exit(1)
	}

	// 程式結束時，結束資料庫連線
	defer func() {
		if err := databaseRepository.DisconnectDB(); err != nil {
			errorLogger.Log("msg", "Disconnect DB fail", "err", err)
			os.Exit(1)
		}
	}()

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		errorLogger.Log("msg", "net listen fail", "err", err)
		os.Exit(1)
	}

	wordService := service.New(logger, databaseRepository)
	wordEndpoints := endpoint.MakeEndpoints(wordService, logger)
	myGrpcServer := transport.NewGRPCServer(wordEndpoints, logger)

	grpcServer := grpc.NewServer()
	pb.RegisterWordServiceServer(grpcServer, myGrpcServer)
	reflection.Register(grpcServer)
	level.Info(logger).Log("msg", "Starting gRPC server at "+PORT)
	err = grpcServer.Serve(listener)

	if err != nil {
		errorLogger.Log("msg", "grpcServer serve fail", "err", err)
	}
}

func loadEnv() {
	switch os.Getenv("WORD_SERVICE_ENV") {
	case "PROD":
		godotenv.Load(".env.production")
	default:
		godotenv.Load(".env.local")
	}
}
