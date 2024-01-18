package config

import (
	"os"
)

func EnvMongoDBURI() string {
	return os.Getenv("MONGODB_URI")
}

func EnvDatabaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func EnvExamServiceServerAddress() string {
	return os.Getenv("EXAM_SERVICE_SERVER_ADDRESS")
}
