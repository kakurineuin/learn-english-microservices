package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/service"
)

type Endpoints struct {
	CreateExam endpoint.Endpoint
	UpdateExam endpoint.Endpoint
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

func MakeEndpoints(examService service.ExamService, logger log.Logger) Endpoints {
	createExamEndpoint := makeCreateExamEndpoint(examService)
	createExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateExam"))(createExamEndpoint)

	updateExamEndpoint := makeUpdateExamEndpoint(examService)
	updateExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "UpdateExam"))(updateExamEndpoint)

	return Endpoints{
		CreateExam: createExamEndpoint,
		UpdateExam: updateExamEndpoint,
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
