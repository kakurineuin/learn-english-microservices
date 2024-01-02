package word

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/wordservice"
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

func NewHandler(wordService wordservice.WordService) WordHandler {
	return &wordHandler{
		wordService: wordService,
	}
}

type wordHandler struct {
	wordService wordservice.WordService
}

func (handler wordHandler) FindWordMeanings(c echo.Context) error {
	errorMessage := "FindWordMeanings failed! error: %w"

	word := c.Param("word")
	userId := utilGetJWTClaims(c).UserId
	c.Logger().Infof("word: %s, userId: %s", word, userId)

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
	var size int32 = 10

	microserviceResponse, err := handler.wordService.FindRandomFavoriteWordMeanings(
		userId,
		size,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}
