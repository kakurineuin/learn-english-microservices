package repository

import (
	"context"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

type transactionFunc func(ctx context.Context) (interface{}, error)

//go:generate mockery --name DatabaseRepository
type DatabaseRepository interface {
	ConnectDB(uri string) error
	DisconnectDB() error

	// Transaction
	WithTransaction(transactoinFunc transactionFunc) (interface{}, error)

	// WordMeaning
	CreateWordMeanings(
		ctx context.Context,
		wordMeanings []model.WordMeaning,
	) (wordMeaningIds []string, err error)
	FindWordMeaningsByWordAndUserId(
		ctx context.Context,
		word, userId string,
	) (wordMeanings []model.WordMeaning, err error)

	// FavoriteWordMeaning
	CreateFavoriteWordMeaning(
		ctx context.Context,
		userId, wordMeaningId string,
	) (favoriteWordMeaningId string, err error)
	GetFavoriteWordMeaningById(
		ctx context.Context,
		favoriteWordMeaningId string,
	) (favoriteWordMeaning *model.FavoriteWordMeaning, err error)
	FindFavoriteWordMeaningsByUserIdAndWord(
		ctx context.Context,
		userId, word string,
		skip, limit int64,
	) (wordMeanings []model.WordMeaning, err error)
	CountFavoriteWordMeaningsByUserIdAndWord(
		ctx context.Context,
		userId, word string,
	) (count int64, err error)
	DeleteFavoriteWordMeaningById(
		ctx context.Context,
		favoriteWordMeaningId string,
	) (deletedCount int64, err error)
}
