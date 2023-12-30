package service

import (
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

type loggingMiddleware struct {
	logger log.Logger
	next   WordService
}

func (mw loggingMiddleware) FindWordByDictionary(
	word, userId string,
) (wordMeanings []model.WordMeaning, err error) {
	defer func() {
		mw.logger.Log("method", "FindWordByDictionary", "word", word, "err", err)
	}()
	return mw.next.FindWordByDictionary(word, userId)
}

func (mw loggingMiddleware) CreateFavoriteWordMeaning(
	userId, wordMeaningId string,
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
	return mw.next.CreateFavoriteWordMeaning(userId, wordMeaningId)
}

func (mw loggingMiddleware) DeleteFavoriteWordMeaning(
	favoriteWordMeaningId, userId string,
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
	return mw.next.DeleteFavoriteWordMeaning(favoriteWordMeaningId, userId)
}

func (mw loggingMiddleware) FindFavoriteWordMeanings(
	pageIndex, pageSize int64,
	userId, word string,
) (total, pageCount int64, favoriteWordMeanings []model.WordMeaning, err error) {
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
	return mw.next.FindFavoriteWordMeanings(pageIndex, pageSize, userId, word)
}
