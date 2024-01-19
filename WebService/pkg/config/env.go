package config

import (
	"os"
)

func EnvJWTSecretKey() string {
	return os.Getenv("JWT_SECRET_KEY")
}

func EnvMongoDBURI() string {
	return os.Getenv("MONGODB_URI")
}

func EnvDatabaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func EnvWebServiceServerAddress() string {
	return os.Getenv("WEB_SERVICE_SERVER_ADDRESS")
}

func EnvExamServiceServerAddress() string {
	return os.Getenv("EXAM_SERVICE_SERVER_ADDRESS")
}

func EnvWordServiceServerAddress() string {
	return os.Getenv("WORD_SERVICE_SERVER_ADDRESS")
}

func EnvEnableCSRF() bool {
	if os.Getenv("ENABLE_CSRF") == "false" {
		return false
	}

	// 若無此環境變數，預設為 true 去啟用 CSRF 檢查
	return true
}
