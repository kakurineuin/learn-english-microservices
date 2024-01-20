package service

import (
	"context"

	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

type loggingMiddleware struct {
	logger log.Logger
	next   WordService
}

func (mw loggingMiddleware) FindWordByDictionary(
	ctx context.Context, word, userId string,
) (wordMeanings []model.WordMeaning, err error) {
	defer func() {
		mw.logger.Log("method", "FindWordByDictionary", "word", word, "err", err)
	}()
	return mw.next.FindWordByDictionary(ctx, word, userId)
}

func (mw loggingMiddleware) CreateFavoriteWordMeaning(
	ctx context.Context, userId, wordMeaningId string,
) (favoriteWordMeaningId string, err error) {
	defer func() {
		mw.logger.Log(
			"method",
			"CreateFavoriteWordMeaning",
			"userId",
			userId,
			"wordMeaningId",
			wordMeaningId,
			"err",
			err,
		)
	}()
	return mw.next.CreateFavoriteWordMeaning(ctx, userId, wordMeaningId)
}

func (mw loggingMiddleware) DeleteFavoriteWordMeaning(
	ctx context.Context, favoriteWordMeaningId, userId string,
) (err error) {
	defer func() {
		mw.logger.Log(
			"method",
			"DeleteFavoriteWordMeaning",
			"favoriteWordMeaningId",
			favoriteWordMeaningId,
			"userId",
			userId,
			"err",
			err,
		)
	}()
	return mw.next.DeleteFavoriteWordMeaning(ctx, favoriteWordMeaningId, userId)
}

func (mw loggingMiddleware) FindFavoriteWordMeanings(
	ctx context.Context,
	pageIndex, pageSize int32,
	userId, word string,
) (total, pageCount int32, favoriteWordMeanings []model.WordMeaning, err error) {
	defer func() {
		mw.logger.Log(
			"method",
			"FindFavoriteWordMeanings",
			"pageIndex",
			pageIndex,
			"pageSize",
			pageSize,
			"userId",
			userId,
			"word",
			word,
			"err",
			err,
		)
	}()
	return mw.next.FindFavoriteWordMeanings(ctx, pageIndex, pageSize, userId, word)
}

func (mw loggingMiddleware) FindRandomFavoriteWordMeanings(
	ctx context.Context, userId string, size int32,
) (wordMeanings []model.WordMeaning, err error) {
	defer func() {
		mw.logger.Log(
			"method",
			"FindRandomFavoriteWordMeanings",
			"userId",
			userId,
			"size",
			size,
			"err",
			err,
		)
	}()
	return mw.next.FindRandomFavoriteWordMeanings(ctx, userId, size)
}
