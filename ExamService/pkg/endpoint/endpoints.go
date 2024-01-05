package endpoint

import (
	"context"
	"time"

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

	CreateQuestion      endpoint.Endpoint
	UpdateQuestion      endpoint.Endpoint
	FindQuestions       endpoint.Endpoint
	DeleteQuestion      endpoint.Endpoint
	FindRandomQuestions endpoint.Endpoint

	CreateExamRecord       endpoint.Endpoint
	FindExamRecords        endpoint.Endpoint
	FindExamRecordOverview endpoint.Endpoint

	FindExamInfos endpoint.Endpoint
}

func MakeEndpoints(examService service.ExamService, logger log.Logger) Endpoints {
	createExamEndpoint := makeCreateExamEndpoint(examService)
	createExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateExam"))(createExamEndpoint)
	createExamEndpoint = RecoverMiddleware(
		log.With(logger, "method", "CreateExam"))(createExamEndpoint)

	updateExamEndpoint := makeUpdateExamEndpoint(examService)
	updateExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "UpdateExam"))(updateExamEndpoint)
	updateExamEndpoint = RecoverMiddleware(
		log.With(logger, "method", "UpdateExam"))(updateExamEndpoint)

	findExamsEndpoint := makeFindExamsEndpoint(examService)
	findExamsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindExams"))(findExamsEndpoint)
	findExamsEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindExams"))(findExamsEndpoint)

	deleteExamEndpoint := makeDeleteExamEndpoint(examService)
	deleteExamEndpoint = LoggingMiddleware(
		log.With(logger, "method", "DeleteExam"))(deleteExamEndpoint)
	deleteExamEndpoint = RecoverMiddleware(
		log.With(logger, "method", "DeleteExam"))(deleteExamEndpoint)

	createQuestionEndpoint := makeCreateQuestionEndpoint(examService)
	createQuestionEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateQuestion"))(createQuestionEndpoint)
	createQuestionEndpoint = RecoverMiddleware(
		log.With(logger, "method", "CreateQuestion"))(createQuestionEndpoint)

	updateQuestionEndpoint := makeUpdateQuestionEndpoint(examService)
	updateQuestionEndpoint = LoggingMiddleware(
		log.With(logger, "method", "UpdateQuestion"))(updateQuestionEndpoint)
	updateQuestionEndpoint = RecoverMiddleware(
		log.With(logger, "method", "UpdateQuestion"))(updateQuestionEndpoint)

	findQuestionsEndpoint := makeFindQuestionsEndpoint(examService)
	findQuestionsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindQuestions"))(findQuestionsEndpoint)
	findQuestionsEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindQuestions"))(findQuestionsEndpoint)

	deleteQuestionEndpoint := makeDeleteQuestionEndpoint(examService)
	deleteQuestionEndpoint = LoggingMiddleware(
		log.With(logger, "method", "DeleteQuestion"))(deleteQuestionEndpoint)
	deleteQuestionEndpoint = RecoverMiddleware(
		log.With(logger, "method", "DeleteQuestion"))(deleteQuestionEndpoint)

	findRandomQuestionsEndpoint := makeFindRandomQuestionsEndpoint(examService)
	findRandomQuestionsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindRandomQuestions"))(findRandomQuestionsEndpoint)
	findRandomQuestionsEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindRandomQuestions"))(findRandomQuestionsEndpoint)

	createExamRecordEndpoint := makeCreateExamRecordEndpoint(examService)
	createExamRecordEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateExamRecord"))(createExamRecordEndpoint)
	createExamRecordEndpoint = RecoverMiddleware(
		log.With(logger, "method", "CreateExamRecord"))(createExamRecordEndpoint)

	findExamRecordsEndpoint := makeFindExamRecordsEndpoint(examService)
	findExamRecordsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindExamRecords"))(findExamRecordsEndpoint)
	findExamRecordsEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindExamRecords"))(findExamRecordsEndpoint)

	findExamRecordOverviewEndpoint := makeFindExamRecordOverviewEndpoint(examService)
	findExamRecordOverviewEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindExamRecordOverview"))(findExamRecordOverviewEndpoint)
	findExamRecordOverviewEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindExamRecordOverview"))(findExamRecordOverviewEndpoint)

	findExamInfosEndpoint := makeFindExamInfosEndpoint(examService)
	findExamInfosEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindExamInfos"))(findExamInfosEndpoint)
	findExamInfosEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindExamInfos"))(findExamInfosEndpoint)

	return Endpoints{
		CreateExam: createExamEndpoint,
		UpdateExam: updateExamEndpoint,
		FindExams:  findExamsEndpoint,
		DeleteExam: deleteExamEndpoint,

		CreateQuestion:      createQuestionEndpoint,
		UpdateQuestion:      updateQuestionEndpoint,
		FindQuestions:       findQuestionsEndpoint,
		DeleteQuestion:      deleteQuestionEndpoint,
		FindRandomQuestions: findRandomQuestionsEndpoint,

		CreateExamRecord:       createExamRecordEndpoint,
		FindExamRecords:        findExamRecordsEndpoint,
		FindExamRecordOverview: findExamRecordOverviewEndpoint,

		FindExamInfos: findExamInfosEndpoint,
	}
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

type FindExamsRequest struct {
	PageIndex int32
	PageSize  int32
	UserId    string
}

type FindExamsResponse struct {
	Total     int32
	PageCount int32
	Exams     []model.Exam
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

type DeleteExamRequest struct {
	ExamId string
	UserId string
}

type DeleteExamResponse struct{}

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

type CreateQuestionRequest struct {
	ExamId  string
	Ask     string
	Answers []string
	UserId  string
}

type CreateQuestionResponse struct {
	QuestionId string
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

type UpdateQuestionRequest struct {
	QuestionId string
	Ask        string
	Answers    []string
	UserId     string
}

type UpdateQuestionResponse struct {
	QuestionId string
}

func makeUpdateQuestionEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateQuestionRequest)
		questionId, err := examService.UpdateQuestion(
			req.QuestionId, req.Ask, req.Answers, req.UserId)
		if err != nil {
			return nil, err
		}
		return UpdateQuestionResponse{QuestionId: questionId}, nil
	}
}

type FindQuestionsRequest struct {
	PageIndex int32
	PageSize  int32
	ExamId    string
	UserId    string
}

type FindQuestionsResponse struct {
	Total     int32
	PageCount int32
	Questions []model.Question
}

func makeFindQuestionsEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindQuestionsRequest)
		total, pageCount, quesitons, err := examService.FindQuestions(
			req.PageIndex,
			req.PageSize,
			req.ExamId,
			req.UserId,
		)
		if err != nil {
			return nil, err
		}
		return FindQuestionsResponse{
			Total:     total,
			PageCount: pageCount,
			Questions: quesitons,
		}, nil
	}
}

type DeleteQuestionRequest struct {
	QuestionId string
	UserId     string
}

type DeleteQuestionResponse struct{}

func makeDeleteQuestionEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteQuestionRequest)
		err := examService.DeleteQuestion(
			req.QuestionId, req.UserId)
		if err != nil {
			return nil, err
		}
		return DeleteQuestionResponse{}, nil
	}
}

type CreateExamRecordRequest struct {
	ExamId           string
	Score            int32
	WrongQuestionIds []string
	UserId           string
}

type CreateExamRecordResponse struct{}

func makeCreateExamRecordEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateExamRecordRequest)
		err := examService.CreateExamRecord(
			req.ExamId, req.Score, req.WrongQuestionIds, req.UserId)
		if err != nil {
			return nil, err
		}
		return CreateExamRecordResponse{}, nil
	}
}

type FindExamRecordsRequest struct {
	PageIndex int32
	PageSize  int32
	ExamId    string
	UserId    string
}

type FindExamRecordsResponse struct {
	Total       int32
	PageCount   int32
	ExamRecords []model.ExamRecord
}

func makeFindExamRecordsEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamRecordsRequest)
		total, pageCount, examRecords, err := examService.FindExamRecords(
			req.PageIndex,
			req.PageSize,
			req.ExamId,
			req.UserId,
		)
		if err != nil {
			return nil, err
		}
		return FindExamRecordsResponse{
			Total:       total,
			PageCount:   pageCount,
			ExamRecords: examRecords,
		}, nil
	}
}

type FindExamRecordOverviewRequest struct {
	ExamId    string
	UserId    string
	StartDate time.Time
}

type FindExamRecordOverviewResponse struct {
	StartDate    string
	Exam         *model.Exam
	Questions    []model.Question
	AnswerWrongs []model.AnswerWrong
	ExamRecords  []model.ExamRecord
}

func makeFindExamRecordOverviewEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamRecordOverviewRequest)
		startDate, exam, quesitons, answerWrongs, examRecords, err := examService.FindExamRecordOverview(
			req.ExamId,
			req.UserId,
			req.StartDate,
		)
		if err != nil {
			return nil, err
		}
		return FindExamRecordOverviewResponse{
			StartDate:    startDate,
			Exam:         exam,
			Questions:    quesitons,
			AnswerWrongs: answerWrongs,
			ExamRecords:  examRecords,
		}, nil
	}
}

type FindExamInfosRequest struct {
	UserId   string
	IsPublic bool
}

type FindExamInfosResponse struct {
	ExamInfos []service.ExamInfo
}

func makeFindExamInfosEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamInfosRequest)
		examInfos, err := examService.FindExamInfos(
			req.UserId,
			req.IsPublic,
		)
		if err != nil {
			return nil, err
		}
		return FindExamInfosResponse{
			ExamInfos: examInfos,
		}, nil
	}
}

type FindRandomQuestionsRequest struct {
	ExamId string
	UserId string
	Size   int32
}

type FindRandomQuestionsResponse struct {
	Exam      *model.Exam
	Questions []model.Question
}

func makeFindRandomQuestionsEndpoint(examService service.ExamService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRandomQuestionsRequest)
		exam, quesitons, err := examService.FindRandomQuestions(
			req.ExamId,
			req.UserId,
			req.Size,
		)
		if err != nil {
			return nil, err
		}
		return FindRandomQuestionsResponse{
			Exam:      exam,
			Questions: quesitons,
		}, nil
	}
}
