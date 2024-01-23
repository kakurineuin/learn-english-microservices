package repository

import (
	"context"
	"time"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
)

type UserHistoryResponse struct {
	Id        string    `json:"_id"       bson:"_id,omitempty"`
	UserId    string    `json:"userId"    bson:"userId"`
	Username  string    `json:"username"  bson:"username"`
	Role      string    `json:"role"      bson:"role"`
	Method    string    `json:"method"    bson:"method"`
	Path      string    `json:"path"      bson:"path"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type transactionFunc func(ctx context.Context) (interface{}, error)

//go:generate mockery --name DatabaseRepository
type DatabaseRepository interface {
	ConnectDB(ctx context.Context, uri string) error
	DisconnectDB(ctx context.Context) error

	// Transaction
	WithTransaction(ctx context.Context, transactoinFunc transactionFunc) (interface{}, error)

	// Exam
	CreateUser(ctx context.Context, user model.User) (userId string, err error)
	GetUserById(ctx context.Context, userId string) (user *model.User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *model.User, err error)
	GetAdminUser(ctx context.Context) (user *model.User, err error)

	// UserHistory
	FindUserHistoryResponsesOrderByUpdatedAt(
		ctx context.Context, pageIndex, pageSize int32,
	) (userHistoryResponses []UserHistoryResponse, err error)
	CountUserHistories(ctx context.Context) (count int32, err error)
	CreateUserHistory(
		ctx context.Context,
		userHistory model.UserHistory,
	) (userHistoryId string, err error)
}

//go:generate mockery --name CacheRepository
type CacheRepository interface {
	ConnectDB(uri string) error
	DisconnectDB() error

	// WordMeaning
	CreateWordMeanings(
		ctx context.Context,
		key string,
		wordMeanings []*pb.WordMeaning,
		expiration time.Duration,
	) error
	FindWordMeanings(ctx context.Context, key string) (wordMeanings []*pb.WordMeaning, err error)
}
