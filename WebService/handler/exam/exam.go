package exam

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kakurineuin/learn-english-microservices/web-service/microservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/util"
)

func CreateExam(c echo.Context) error {
	type RequestBody struct {
		Topic       string `json:"topic"`
		Description string `json:"description"`
		IsPublic    bool   `json:"isPublic"`
	}

	errorMessage := "CreateExam failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := util.GetJWTClaims(c).UserId

	microserviceResponse, err := microservice.CreateExam(
		requestBody.Topic,
		requestBody.Description,
		requestBody.IsPublic,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"examId": microserviceResponse.ExamId,
	})
}

func FindExams(c echo.Context) error {
	errorMessage := "FindExams failed! error: %w"

	var (
		pageIndex int64 = 0
		pageSize  int64 = 0
	)

	err := echo.QueryParamsBinder(c).
		Int64("pageIndex", &pageIndex).
		Int64("pageSize", &pageSize).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := util.GetJWTClaims(c).UserId

	microserviceResponse, err := microservice.FindExams(
		pageIndex,
		pageSize,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	result, err := protojson.MarshalOptions{
		EmitUnpopulated: true, // Zero value 的欄位不要省略
	}.Marshal(microserviceResponse)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.JSONBlob(http.StatusOK, result)
}

func UpdateExam(c echo.Context) error {
	type RequestBody struct {
		ExamId      string `json:"_id"`
		Topic       string `json:"topic"`
		Description string `json:"description"`
		IsPublic    bool   `json:"isPublic"`
	}

	errorMessage := "UpdateExam failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := util.GetJWTClaims(c).UserId

	microserviceResponse, err := microservice.UpdateExam(
		requestBody.ExamId,
		requestBody.Topic,
		requestBody.Description,
		requestBody.IsPublic,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"examId": microserviceResponse.ExamId,
	})
}
