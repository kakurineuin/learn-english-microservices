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

	return Endpoints{
		CreateExam: createExamEndpoint,
		UpdateExam: updateExamEndpoint,
		FindExams:  findExamsEndpoint,
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
