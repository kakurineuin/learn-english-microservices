package endpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

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
	var createExamEndpoint endpoint.Endpoint
	{
		createExamEndpoint = makeCreateExamEndpoint(examService)
		createExamEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			createExamEndpoint,
		)
		createExamEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			createExamEndpoint,
		)
		createExamEndpoint = LoggingMiddleware(
			log.With(logger, "method", "CreateExam"))(createExamEndpoint)
		createExamEndpoint = RecoverMiddleware(
			log.With(logger, "method", "CreateExam"))(createExamEndpoint)
	}

	var updateExamEndpoint endpoint.Endpoint
	{
		updateExamEndpoint = makeUpdateExamEndpoint(examService)
		updateExamEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			updateExamEndpoint,
		)
		updateExamEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			updateExamEndpoint,
		)
		updateExamEndpoint = LoggingMiddleware(
			log.With(logger, "method", "UpdateExam"))(updateExamEndpoint)
		updateExamEndpoint = RecoverMiddleware(
			log.With(logger, "method", "UpdateExam"))(updateExamEndpoint)
	}

	var findExamsEndpoint endpoint.Endpoint
	{
		findExamsEndpoint = makeFindExamsEndpoint(examService)
		findExamsEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findExamsEndpoint,
		)
		findExamsEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findExamsEndpoint,
		)
		findExamsEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindExams"))(findExamsEndpoint)
		findExamsEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindExams"))(findExamsEndpoint)
	}

	var deleteExamEndpoint endpoint.Endpoint
	{
		deleteExamEndpoint = makeDeleteExamEndpoint(examService)
		deleteExamEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			deleteExamEndpoint,
		)
		deleteExamEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			deleteExamEndpoint,
		)
		deleteExamEndpoint = LoggingMiddleware(
			log.With(logger, "method", "DeleteExam"))(deleteExamEndpoint)
		deleteExamEndpoint = RecoverMiddleware(
			log.With(logger, "method", "DeleteExam"))(deleteExamEndpoint)
	}

	var createQuestionEndpoint endpoint.Endpoint
	{
		createQuestionEndpoint = makeCreateQuestionEndpoint(examService)
		createQuestionEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			createQuestionEndpoint,
		)
		createQuestionEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			createQuestionEndpoint,
		)
		createQuestionEndpoint = LoggingMiddleware(
			log.With(logger, "method", "CreateQuestion"))(createQuestionEndpoint)
		createQuestionEndpoint = RecoverMiddleware(
			log.With(logger, "method", "CreateQuestion"))(createQuestionEndpoint)
	}

	var updateQuestionEndpoint endpoint.Endpoint
	{
		updateQuestionEndpoint = makeUpdateQuestionEndpoint(examService)
		updateQuestionEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			updateQuestionEndpoint,
		)
		updateQuestionEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			updateQuestionEndpoint,
		)
		updateQuestionEndpoint = LoggingMiddleware(
			log.With(logger, "method", "UpdateQuestion"))(updateQuestionEndpoint)
		updateQuestionEndpoint = RecoverMiddleware(
			log.With(logger, "method", "UpdateQuestion"))(updateQuestionEndpoint)
	}

	var findQuestionsEndpoint endpoint.Endpoint
	{
		findQuestionsEndpoint = makeFindQuestionsEndpoint(examService)
		findQuestionsEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findQuestionsEndpoint,
		)
		findQuestionsEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findQuestionsEndpoint,
		)
		findQuestionsEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindQuestions"))(findQuestionsEndpoint)
		findQuestionsEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindQuestions"))(findQuestionsEndpoint)
	}

	var deleteQuestionEndpoint endpoint.Endpoint
	{
		deleteQuestionEndpoint = makeDeleteQuestionEndpoint(examService)
		deleteQuestionEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			deleteQuestionEndpoint,
		)
		deleteQuestionEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			deleteQuestionEndpoint,
		)
		deleteQuestionEndpoint = LoggingMiddleware(
			log.With(logger, "method", "DeleteQuestion"))(deleteQuestionEndpoint)
		deleteQuestionEndpoint = RecoverMiddleware(
			log.With(logger, "method", "DeleteQuestion"))(deleteQuestionEndpoint)
	}

	var findRandomQuestionsEndpoint endpoint.Endpoint
	{
		findRandomQuestionsEndpoint = makeFindRandomQuestionsEndpoint(examService)
		findRandomQuestionsEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findRandomQuestionsEndpoint,
		)
		findRandomQuestionsEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findRandomQuestionsEndpoint,
		)
		findRandomQuestionsEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindRandomQuestions"))(findRandomQuestionsEndpoint)
		findRandomQuestionsEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindRandomQuestions"))(findRandomQuestionsEndpoint)
	}

	var createExamRecordEndpoint endpoint.Endpoint
	{
		createExamRecordEndpoint = makeCreateExamRecordEndpoint(examService)
		createExamRecordEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			createExamRecordEndpoint,
		)
		createExamRecordEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			createExamRecordEndpoint,
		)
		createExamRecordEndpoint = LoggingMiddleware(
			log.With(logger, "method", "CreateExamRecord"))(createExamRecordEndpoint)
		createExamRecordEndpoint = RecoverMiddleware(
			log.With(logger, "method", "CreateExamRecord"))(createExamRecordEndpoint)
	}

	var findExamRecordsEndpoint endpoint.Endpoint
	{
		findExamRecordsEndpoint = makeFindExamRecordsEndpoint(examService)
		findExamRecordsEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findExamRecordsEndpoint,
		)
		findExamRecordsEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findExamRecordsEndpoint,
		)
		findExamRecordsEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindExamRecords"))(findExamRecordsEndpoint)
		findExamRecordsEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindExamRecords"))(findExamRecordsEndpoint)
	}

	var findExamRecordOverviewEndpoint endpoint.Endpoint
	{
		findExamRecordOverviewEndpoint = makeFindExamRecordOverviewEndpoint(examService)
		findExamRecordOverviewEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findExamRecordOverviewEndpoint,
		)
		findExamRecordOverviewEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findExamRecordOverviewEndpoint,
		)
		findExamRecordOverviewEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindExamRecordOverview"))(findExamRecordOverviewEndpoint)
		findExamRecordOverviewEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindExamRecordOverview"))(findExamRecordOverviewEndpoint)
	}

	var findExamInfosEndpoint endpoint.Endpoint
	{
		findExamInfosEndpoint = makeFindExamInfosEndpoint(examService)
		findExamInfosEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findExamInfosEndpoint,
		)
		findExamInfosEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findExamInfosEndpoint,
		)
		findExamInfosEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindExamInfos"))(findExamInfosEndpoint)
		findExamInfosEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindExamInfos"))(findExamInfosEndpoint)
	}

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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateExamRequest)
		examId, err := examService.CreateExam(
			ctx, req.Topic, req.Description, req.IsPublic, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateExamRequest)
		examId, err := examService.UpdateExam(
			ctx, req.ExamId, req.Topic, req.Description, req.IsPublic, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamsRequest)
		total, pageCount, exams, err := examService.FindExams(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteExamRequest)
		err := examService.DeleteExam(ctx, req.ExamId, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateQuestionRequest)
		questionId, err := examService.CreateQuestion(
			ctx, req.ExamId, req.Ask, req.Answers, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateQuestionRequest)
		questionId, err := examService.UpdateQuestion(
			ctx, req.QuestionId, req.Ask, req.Answers, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindQuestionsRequest)
		total, pageCount, quesitons, err := examService.FindQuestions(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteQuestionRequest)
		err := examService.DeleteQuestion(
			ctx, req.QuestionId, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateExamRecordRequest)
		err := examService.CreateExamRecord(
			ctx, req.ExamId, req.Score, req.WrongQuestionIds, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamRecordsRequest)
		total, pageCount, examRecords, err := examService.FindExamRecords(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamRecordOverviewRequest)
		startDate, exam, quesitons, answerWrongs, examRecords, err := examService.FindExamRecordOverview(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindExamInfosRequest)
		examInfos, err := examService.FindExamInfos(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRandomQuestionsRequest)
		exam, quesitons, err := examService.FindRandomQuestions(
			ctx,
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
