package word

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kakurineuin/learn-english-microservices/web-service/microservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/util"
)

func FindWordMeanings(c echo.Context) error {
	errorMessage := "FindWordMeanings failed! error: %w"

	word := c.Param("word")
	userId := util.GetJWTClaims(c).UserId
	c.Logger().Infof("word: %s, userId: %s", word, userId)

	microserviceResponse, err := microservice.FindWordByDictionary(word, userId)
	if err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusInternalServerError)
	}

	result, err := protojson.MarshalOptions{
		EmitUnpopulated: true, // Zero value 的欄位不要省略
	}.Marshal(microserviceResponse)
	if err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, result)
}

func CreateFavoriteWordMeaning(c echo.Context) error {
	type RequestBody struct {
		WordMeaningId string `json:"wordMeaningId"`
	}

	errorMessage := "CreateFavoriteWordMeaning failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusBadRequest)
	}

	userId := util.GetJWTClaims(c).UserId
	c.Logger().Infof("requestBody: %v, userId: %s", requestBody, userId)

	microserviceResponse, err := microservice.CreateFavoriteWordMeaning(
		userId,
		requestBody.WordMeaningId,
	)
	if err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"favoriteWordMeaningId": microserviceResponse.FavoriteWordMeaningId,
	})
}

func DeleteFavoriteWordMeaning(c echo.Context) error {
	errorMessage := "DeleteFavoriteWordMeaning failed! error: %w"
	favoriteWordMeaningId := c.Param("favoriteWordMeaningId")
	userId := util.GetJWTClaims(c).UserId
	c.Logger().Infof("favoriteWordMeaningId: %s, userId: %s", favoriteWordMeaningId, userId)

	_, err := microservice.DeleteFavoriteWordMeaning(
		favoriteWordMeaningId,
		userId,
	)
	if err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func FindFavoriteWordMeanings(c echo.Context) error {
	errorMessage := "FindFavoriteWordMeanings failed! error: %w"

	var (
		pageIndex int64  = 0
		pageSize  int64  = 0
		word      string = ""
	)

	err := echo.QueryParamsBinder(c).
		Int64("pageIndex", &pageIndex).
		Int64("pageSize", &pageSize).
		String("word", &word).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusBadRequest)
	}

	userId := util.GetJWTClaims(c).UserId

	microserviceResponse, err := microservice.FindFavoriteWordMeanings(
		pageIndex,
		pageSize,
		userId,
		word,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONError(c, http.StatusInternalServerError)
	}

	result, err := protojson.MarshalOptions{
		EmitUnpopulated: true, // Zero value 的欄位不要省略
	}.Marshal(microserviceResponse)
	if err != nil {
		c.Logger().Errorf(errorMessage, err)
		return util.SendJSONError(c, http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, result)
}
