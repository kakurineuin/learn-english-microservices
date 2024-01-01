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

func EnvServerAddress() string {
	return os.Getenv("SERVER_ADDRESS")
}

func EnvExamServiceServerAddress() string {
	return os.Getenv("EXAM_SERVICE_SERVER_ADDRESS")
}

func EnvWordServiceServerAddress() string {
	return os.Getenv("WORD_SERVICE_SERVER_ADDRESS")
}
