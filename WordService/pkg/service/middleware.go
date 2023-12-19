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
