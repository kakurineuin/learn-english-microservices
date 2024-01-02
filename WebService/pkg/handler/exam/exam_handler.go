package exam

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/examservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

// For mock at test
var utilGetJWTClaims = util.GetJWTClaims

type ExamHandler interface {
	CreateExam(c echo.Context) error
	FindExams(c echo.Context) error
	UpdateExam(c echo.Context) error
	DeleteExam(c echo.Context) error
	FindQuestions(c echo.Context) error
	CreateQuestion(c echo.Context) error
	UpdateQuestion(c echo.Context) error
	DeleteQuestion(c echo.Context) error
	FindRandomQuestions(c echo.Context) error
}

type examHandler struct {
	examServce examservice.ExamService
}

func NewHandler(examServce examservice.ExamService) ExamHandler {
	return &examHandler{
		examServce: examServce,
	}
}

func (handler examHandler) CreateExam(c echo.Context) error {
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

	userId := utilGetJWTClaims(c).UserId

	microserviceResponse, err := handler.examServce.CreateExam(
		requestBody.Topic,
		requestBody.Description,
		requestBody.IsPublic,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler examHandler) FindExams(c echo.Context) error {
	errorMessage := "FindExams failed! error: %w"

	var (
		pageIndex int32 = 0
		pageSize  int32 = 0
	)

	err := echo.QueryParamsBinder(c).
		Int32("pageIndex", &pageIndex).
		Int32("pageSize", &pageSize).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId

	microserviceResponse, err := handler.examServce.FindExams(
		pageIndex,
		pageSize,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler examHandler) UpdateExam(c echo.Context) error {
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

	userId := utilGetJWTClaims(c).UserId

	microserviceResponse, err := handler.examServce.UpdateExam(
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

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler examHandler) DeleteExam(c echo.Context) error {
	errorMessage := "DeleteExam failed! error: %w"

	examId := ""
	err := echo.PathParamsBinder(c).
		String("examId", &examId).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId

	_, err = handler.examServce.DeleteExam(examId, userId)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.NoContent(http.StatusOK)
}

func (handler examHandler) FindQuestions(c echo.Context) error {
	errorMessage := "FindQuestions failed! error: %w"

	var (
		examId    string = ""
		pageIndex int32  = 0
		pageSize  int32  = 0
	)

	err := echo.PathParamsBinder(c).
		String("examId", &examId).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	err = echo.QueryParamsBinder(c).
		Int32("pageIndex", &pageIndex).
		Int32("pageSize", &pageSize).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId

	microserviceResponse, err := handler.examServce.FindQuestions(
		pageIndex,
		pageSize,
		examId,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler examHandler) CreateQuestion(c echo.Context) error {
	type RequestBody struct {
		Ask     string   `json:"ask"`
		Answers []string `json:"answers"`
	}

	errorMessage := "CreateQuestion failed! error: %w"

	examId := ""
	err := echo.PathParamsBinder(c).
		String("examId", &examId).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId
	microserviceResponse, err := handler.examServce.CreateQuestion(
		examId,
		requestBody.Ask,
		requestBody.Answers,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler examHandler) UpdateQuestion(c echo.Context) error {
	type RequestBody struct {
		QuestionId string   `json:"_id"`
		Ask        string   `json:"ask"`
		Answers    []string `json:"answers"`
	}

	errorMessage := "UpdateQuestion failed! error: %w"

	requestBody := new(RequestBody)
	if err := c.Bind(&requestBody); err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId
	microserviceResponse, err := handler.examServce.UpdateQuestion(
		requestBody.QuestionId,
		requestBody.Ask,
		requestBody.Answers,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}

func (handler examHandler) DeleteQuestion(c echo.Context) error {
	errorMessage := "DeleteQuestion failed! error: %w"

	questionId := ""
	err := echo.PathParamsBinder(c).
		String("questionId", &questionId).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId
	_, err = handler.examServce.DeleteQuestion(
		questionId,
		userId,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return c.NoContent(http.StatusOK)
}

func (handler examHandler) FindRandomQuestions(c echo.Context) error {
	errorMessage := "FindRandomQuestions failed! error: %w"

	examId := ""
	err := echo.PathParamsBinder(c).
		String("examId", &examId).
		BindError() // returns first binding error
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONBadRequest(c)
	}

	userId := utilGetJWTClaims(c).UserId
	var size int32 = 10

	microserviceResponse, err := handler.examServce.FindRandomQuestions(
		examId,
		userId,
		size,
	)
	if err != nil {
		c.Logger().Error(fmt.Errorf(errorMessage, err))
		return util.SendJSONInternalServerError(c)
	}

	return util.SendJSONResponse(c, microserviceResponse)
}
