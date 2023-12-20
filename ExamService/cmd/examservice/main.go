package main

import (
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/joho/godotenv"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pb"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/database"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/endpoint"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/service"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/transport"
)

const PORT = ":8090"

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	errorLogger := level.Error(logger)

	loadEnv()

	// 連線到資料庫
	if err := database.ConnectDB(); err != nil {
		errorLogger.Log("msg", "Connect DB fail", "err", err)
		os.Exit(1)
	}

	// 程式結束時，結束資料庫連線
	defer func() {
		if err := database.DisconnectDB(); err != nil {
			errorLogger.Log("msg", "Disconnect DB fail", "err", err)
			os.Exit(1)
		}
	}()

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		errorLogger.Log("msg", "net listen fail", "err", err)
		os.Exit(1)
	}

	examService := service.New(logger)
	examEndpoints := endpoint.MakeEndpoints(examService, logger)
	myGrpcServer := transport.NewGRPCServer(examEndpoints, logger)

	grpcServer := grpc.NewServer()
	pb.RegisterExamServiceServer(grpcServer, myGrpcServer)
	reflection.Register(grpcServer)
	level.Info(logger).Log("msg", "Starting gRPC server at "+PORT)
	err = grpcServer.Serve(listener)

	if err != nil {
		errorLogger.Log("msg", "grpcServer serve fail", "err", err)
	}
}

func loadEnv() {
	switch os.Getenv("EXAM_SERVICE_ENV") {
	case "PROD":
		godotenv.Load(".env.production")
	default:
		godotenv.Load(".env.local")
	}
}