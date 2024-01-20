package config

import (
	"os"
)

func EnvEnableTransaction() bool {
	// 預設啟用 MongoDB 交易功能
	return os.Getenv("ENABLE_TRANSACTION") != "false"
}

func EnvMongoDBURI() string {
	return os.Getenv("MONGODB_URI")
}

func EnvDatabaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func EnvExamServiceServerAddress() string {
	return os.Getenv("EXAM_SERVICE_SERVER_ADDRESS")
}
