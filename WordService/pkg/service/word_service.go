package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/crawler"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/repository"
)

type WordService interface {
	FindWordByDictionary(word, userId string) ([]model.WordMeaning, error)
	CreateFavoriteWordMeaning(
		userId, wordMeaningId string,
	) (favoriteWordMeaningId string, err error)
}

type wordService struct {
	logger             log.Logger
	errorLogger        log.Logger
	databaseRepository repository.DatabaseRepository
	spider             crawler.Spider
}

func New(logger log.Logger, databaseRepository repository.DatabaseRepository) WordService {
	var wordService WordService = wordService{
		logger:             logger,
		errorLogger:        level.Error(logger),
		databaseRepository: databaseRepository,
		spider:             crawler.NewSpider(),
	}
	wordService = loggingMiddleware{logger, wordService}
	return wordService
}

func (wordService wordService) FindWordByDictionary(
	word, userId string,
) ([]model.WordMeaning, error) {
	logger := wordService.logger
	errorLogger := wordService.errorLogger
	errorMessage := "FindWordByDictionary failed! error: %w"
	logger.Log("msg", "Start FindWordByDictionary", "word", word, "userId", userId)

	databaseRepository := wordService.databaseRepository
	spider := wordService.spider

	// 統一以小寫去查詢
	word = strings.ToLower(word)
	wordMeanings, err := databaseRepository.FindWordMeaningsByWordAndUserId(
		context.TODO(),
		word,
		userId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	// 若資料庫尚無此單字的資料
	if len(wordMeanings) == 0 {
		logger.Log("msg", fmt.Sprintf("Start crawling dictionary website by word: %s", word))
		wordMeanings, err = spider.FindWordMeaningsFromDictionary(word)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, fmt.Errorf(errorMessage, err)
		}

		logger.Log(
			"msg",
			fmt.Sprintf("Crawling dictionary website result size: %d", len(wordMeanings)),
		)

		// 如果線上辭典網站查無此單字的解釋，那就是查無資料
		if len(wordMeanings) == 0 {
			return nil, nil
		}

		// 新增到資料庫
		_, err = wordService.databaseRepository.CreateWordMeanings(context.TODO(), wordMeanings)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, fmt.Errorf(errorMessage, err)
		}

		// 從資料庫查詢後再回傳，這樣每筆資料就會有正確的 mongodb _id
		wordMeanings, err = databaseRepository.FindWordMeaningsByWordAndUserId(
			context.TODO(),
			word,
			userId,
		)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, fmt.Errorf(errorMessage, err)
		}
	}

	return wordMeanings, nil
}

func (wordService wordService) CreateFavoriteWordMeaning(
	userId, wordMeaningId string,
) (favoriteWordMeaningId string, err error) {
	errorLogger := wordService.errorLogger
	errorMessage := "CreateFavoriteWordMeaning failed! error: %w"

	databaseRepository := wordService.databaseRepository
	favoriteWordMeaningId, err = databaseRepository.CreateFavoriteWordMeaning(
		context.TODO(),
		userId,
		wordMeaningId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	return favoriteWordMeaningId, nil
}
