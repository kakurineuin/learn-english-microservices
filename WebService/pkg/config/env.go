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
