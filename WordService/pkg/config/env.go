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

func EnvServerAddress() string {
	return os.Getenv("SERVER_ADDRESS")
}
