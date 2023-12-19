package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/service"
)

type Endpoints struct {
	CreateExam endpoint.Endpoint
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

func MakeEndpoints(examService service.ExamService, logger log.Logger) Endpoints {
	createExamEndpoint := makeCreateExamEndpoint(examService)
	createExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateExam"))(createExamEndpoint)

	return Endpoints{
		CreateExam: createExamEndpoint,
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
