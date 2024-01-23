package word

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/wordservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

// For mock at test
var utilGetJWTClaims = util.GetJWTClaims

type WordHandler interface {
	FindWordMeanings(c echo.Context) error
	CreateFavoriteWordMeaning(c echo.Context) error
	DeleteFavoriteWordMeaning(c echo.Context) error
	FindFavoriteWordMeanings(c echo.Context) error
	FindRandomFavoriteWordMeanings(c echo.Context) error
}

func NewHandler(
	wordService wordservice.WordService,
	cacheRepository repository.CacheRepository,
) WordHandler {
	return &wordHandler{
		wordService:     wordService,
		cacheRepository: cacheRepository,
	}
}

type wordHandler struct {
	wordService     wordservice.WordService
	cacheRepository repository.CacheRepository
}

func (handler wordHandler) FindWordMeanings(c echo.Context) error {
	errorMessage := "FindWordMeanings failed! error: %w"

	word := c.Param("word")
	userId := utilGetJWTClaims(c).UserId
	c.Logger().Infof("============== word: %s, userId: %s", word, userId)

	microserviceResponse, err := handler.wordService.FindWordByDictionary(word, userId)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler wordHandler) CreateFavoriteWordMeaning(c echo.Context) error {
	type RequestBody struct {
		WordMeaningId string `json:"wordMeaningId"`
	}

	errorMessage := "CreateFavoriteWordMeaning failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId
	c.Logger().Infof("requestBody: %v, userId: %s", requestBody, userId)

	microserviceResponse, err := handler.wordService.CreateFavoriteWordMeaning(
		userId,
		requestBody.WordMeaningId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler wordHandler) DeleteFavoriteWordMeaning(c echo.Context) error {
	errorMessage := "DeleteFavoriteWordMeaning failed! error: %w"
	favoriteWordMeaningId := c.Param("favoriteWordMeaningId")
	userId := utilGetJWTClaims(c).UserId
	c.Logger().Infof("favoriteWordMeaningId: %s, userId: %s", favoriteWordMeaningId, userId)

	_, err := handler.wordService.DeleteFavoriteWordMeaning(
		favoriteWordMeaningId,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.NoContent(http.StatusOK)
}

func (handler wordHandler) FindFavoriteWordMeanings(c echo.Context) error {
	errorMessage := "FindFavoriteWordMeanings failed! error: %w"

	var (
		pageIndex int32  = 0
		pageSize  int32  = 0
		word      string = ""
	)

	err := echo.QueryParamsBinder(c).
		Int32("pageIndex", &pageIndex).
		Int32("pageSize", &pageSize).
		String("word", &word).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId

	microserviceResponse, err := handler.wordService.FindFavoriteWordMeanings(
		pageIndex,
		pageSize,
		userId,
		word,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler wordHandler) FindRandomFavoriteWordMeanings(c echo.Context) error {
	errorMessage := "FindRandomFavoriteWordMeanings failed! error: %w"

	userId := utilGetJWTClaims(c).UserId
	returnSize := 10

	// 先檢查 cache
	cacheRepository := handler.cacheRepository
	key := fmt.Sprintf("FindRandomFavoriteWordMeanings:%s", userId)
	wordMeanings, err := cacheRepository.FindWordMeanings(c.Request().Context(), key)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	// 若 cache 無資料，就呼叫 WordService
	if len(wordMeanings) == 0 {
		var maxSize int32 = 20
		microserviceResponse, err := handler.wordService.FindRandomFavoriteWordMeanings(
			userId,
			maxSize,
		)
		if err != nil {
			c.Logger().Error(fmt.Errorf(errorMessage, err))
			return util.SendJSONInternalServerError(c)
		}

		wordMeanings = microserviceResponse.FavoriteWordMeanings

		// 資料夠多就放入 cache
		if len(wordMeanings) >= returnSize {
			err = cacheRepository.CreateWordMeanings(
				c.Request().Context(),
				key,
				wordMeanings,
				5*time.Minute,
			)
			if err != nil {
				c.Logger().Error(fmt.Errorf(errorMessage, err))
				return util.SendJSONInternalServerError(c)
			}
		}
	}

	wordMeaningsSize := len(wordMeanings)

	if wordMeaningsSize > 1 {

		// 隨機洗牌變更單字解釋的順序
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		rng.Shuffle(wordMeaningsSize, func(i, j int) {
			wordMeanings[i], wordMeanings[j] = wordMeanings[j], wordMeanings[i]
		})
	}

	if wordMeaningsSize > returnSize {
		wordMeanings = wordMeanings[0:returnSize]
	}

	return util.SendJSONResponse(c, &pb.FindRandomFavoriteWordMeaningsResponse{
		FavoriteWordMeanings: wordMeanings,
	})
}
