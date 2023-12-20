package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/service"
)

type Endpoints struct {
	CreateExam endpoint.Endpoint
	UpdateExam endpoint.Endpoint
	FindExams  endpoint.Endpoint
	DeleteExam endpoint.Endpoint

	CreateQuestion endpoint.Endpoint
}

type CreateExamRequest struct {
	Topic       string
	Description string
	IsPublic    bool
	UserId      string
}

type CreateExamResponse struct {
	ExamId string
}

type UpdateExamRequest struct {
	ExamId      string
	Topic       string
	Description string
	IsPublic    bool
	UserId      string
}

type UpdateExamResponse struct {
	ExamId string
}

type FindExamsRequest struct {
	PageIndex int64
	PageSize  int64
	UserId    string
}

type FindExamsResponse struct {
	Total     int64
	PageCount int64
	Exams     []model.Exam
}

type DeleteExamRequest struct {
	ExamId string
	UserId string
}

type DeleteExamResponse struct{}

type CreateQuestionRequest struct {
	ExamId  string
	Ask     string
	Answers []string
	UserId  string
}

type CreateQuestionResponse struct {
	QuestionId string
}

func MakeEndpoints(examService service.ExamService, logger log.Logger) Endpoints {
	createExamEndpoint := makeCreateExamEndpoint(examService)
	createExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateExam"))(createExamEndpoint)

	updateExamEndpoint := makeUpdateExamEndpoint(examService)
	updateExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "UpdateExam"))(updateExamEndpoint)

	findExamsEndpoint := makeFindExamsEndpoint(examService)
	findExamsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindExams"))(findExamsEndpoint)

	deleteExamEndpoint := makeDeleteExamEndpoint(examService)
	deleteExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "DeleteExam"))(deleteExamEndpoint)

	createQuestionEndpoint := makeCreateQuestionEndpoint(examService)
	createQuestionEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateQuestion"))(createQuestionEndpoint)

	return Endpoints{
		CreateExam: createExamEndpoint,
		UpdateExam: updateExamEndpoint,
		FindExams:  findExamsEndpoint,
		DeleteExam: deleteExamEndpoint,

		CreateQuestion: createQuestionEndpoint,
	}
}

func makeCreateExamEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateExamRequest)
		examId, err := examService.CreateExam(
			req.Topic, req.Description, req.IsPublic, req.UserId)
		if err != nil {
			return nil, err
		}
		return CreateExamResponse{ExamId: examId}, nil
	}
}

func makeUpdateExamEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateExamRequest)
		examId, err := examService.UpdateExam(
			req.ExamId, req.Topic, req.Description, req.IsPublic, req.UserId)
		if err != nil {
			return nil, err
		}
		return UpdateExamResponse{ExamId: examId}, nil
	}
}

func makeFindExamsEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamsRequest)
		total, pageCount, exams, err := examService.FindExams(
			req.PageIndex,
			req.PageSize,
			req.UserId,
		)
		if err != nil {
			return nil, err
		}
		return FindExamsResponse{
			Total:     total,
			PageCount: pageCount,
			Exams:     exams,
		}, nil
	}
}

func makeDeleteExamEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteExamRequest)
		err := examService.DeleteExam(req.ExamId, req.UserId)
		if err != nil {
			return nil, err
		}
		return DeleteExamResponse{}, nil
	}
}

func makeCreateQuestionEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateQuestionRequest)
		questionId, err := examService.CreateQuestion(
			req.ExamId, req.Ask, req.Answers, req.UserId)
		if err != nil {
			return nil, err
		}
		return CreateQuestionResponse{QuestionId: questionId}, nil
	}
}
