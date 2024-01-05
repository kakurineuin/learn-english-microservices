package repository

import (
	"context"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
)

type transactionFunc func(ctx context.Context) (interface{}, error)

//go:generate mockery --name DatabaseRepository
type DatabaseRepository interface {
	ConnectDB(uri string) error
	DisconnectDB() error

	// Transaction
	WithTransaction(transactoinFunc transactionFunc) (interface{}, error)

	// Exam
	CreateUser(ctx context.Context, user model.User) (userId string, err error)
	GetUserById(ctx context.Context, userId string) (user *model.User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *model.User, err error)
	GetAdminUser(ctx context.Context) (user *model.User, err error)
}
