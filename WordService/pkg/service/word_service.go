package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/crawler"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/repository"
)

var unauthorizedOperationError = fmt.Errorf("Unauthorized operation")

type WordService interface {
	FindWordByDictionary(word, userId string) ([]model.WordMeaning, error)
	CreateFavoriteWordMeaning(
		userId, wordMeaningId string,
	) (favoriteWordMeaningId string, err error)
	DeleteFavoriteWordMeaning(
		favoriteWordMeaningId, userId string,
	) error
	FindFavoriteWordMeanings(
		pageIndex, pageSize int32,
		userId, word string,
	) (total, pageCount int32, favoriteWordMeanings []model.WordMeaning, err error)
	FindRandomFavoriteWordMeanings(
		userId string, size int32,
	) (wordMeanings []model.WordMeaning, err error)
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

func (wordService wordService) DeleteFavoriteWordMeaning(
	favoriteWordMeaningId, userId string,
) error {
	errorLogger := wordService.errorLogger
	errorMessage := "DeleteFavoriteWordMeaning failed! error: %w"

	databaseRepository := wordService.databaseRepository

	favoriteWordMeaning, err := databaseRepository.GetFavoriteWordMeaningById(
		context.TODO(),
		favoriteWordMeaningId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	if favoriteWordMeaning == nil {
		err = fmt.Errorf("FavoriteWordMeaning not found by id: %s", favoriteWordMeaningId)
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	// 檢查不能刪除別人的資料
	if favoriteWordMeaning.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	_, err = databaseRepository.DeleteFavoriteWordMeaningById(
		context.TODO(),
		favoriteWordMeaningId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	return nil
}

func (wordService wordService) FindFavoriteWordMeanings(
	pageIndex, pageSize int32,
	userId, word string,
) (total, pageCount int32, favoriteWordMeanings []model.WordMeaning, err error) {
	errorLogger := wordService.errorLogger
	errorMessage := "FindFavoriteWordMeanings failed! error: %w"

	skip := pageSize * pageIndex
	limit := pageSize
	databaseRepository := wordService.databaseRepository
	favoriteWordMeanings, err = databaseRepository.FindFavoriteWordMeaningsByUserIdAndWord(
		context.TODO(),
		userId,
		word,
		skip,
		limit,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	total, err = databaseRepository.CountFavoriteWordMeaningsByUserIdAndWord(
		context.TODO(),
		userId,
		word,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// PageCount
	pageCount = int32(math.Ceil(float64(total) / float64(pageSize)))
	return total, pageCount, favoriteWordMeanings, nil
}

func (wordService wordService) FindRandomFavoriteWordMeanings(
	userId string, size int32,
) (wordMeanings []model.WordMeaning, err error) {
	errorLogger := wordService.errorLogger
	errorMessage := "FindRandomFavoriteWordMeanings failed! error: %w"

	databaseRepository := wordService.databaseRepository
	total, err := databaseRepository.CountFavoriteWordMeaningsByUserIdAndWord(
		context.TODO(),
		userId,
		"",
	)
	if err != nil {
		errorLogger.Log("err", err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	// 總數是零，表示使用者還沒新增喜歡的單字解釋
	if total == 0 {
		return []model.WordMeaning{}, nil
	}

	indexes := []int32{}

	for i := int32(0); i < total; i++ {
		indexes = append(indexes, i)
	}

	// 將 indexes 隨機洗牌，然後根據洗牌後的 indexes 順序去查詢，達到隨機排序的效果
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})

	maxQueryTotal := min(total, size)

	for i := int32(0); i < maxQueryTotal; i++ {
		findWordMeanings, err := databaseRepository.FindFavoriteWordMeaningsByUserIdAndWord(
			context.TODO(),
			userId,
			"",
			indexes[i],
			1,
		)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, fmt.Errorf(errorMessage, err)
		}

		wordMeanings = append(wordMeanings, findWordMeanings[0])
	}

	return wordMeanings, nil
}
