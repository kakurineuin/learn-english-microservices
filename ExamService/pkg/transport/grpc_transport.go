package transport

import (
	"context"
	"errors"
	"time"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pb"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/endpoint"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

const TIMEOUT = 20 * time.Second

type GRPCServer struct {
	logger log.Logger

	createExam gt.Handler
	updateExam gt.Handler
	findExams  gt.Handler
	deleteExam gt.Handler

	createQuestion      gt.Handler
	updateQuestion      gt.Handler
	findQuestions       gt.Handler
	deleteQuestion      gt.Handler
	findRandomQuestions gt.Handler

	createExamRecord       gt.Handler
	findExamRecords        gt.Handler
	findExamRecordOverview gt.Handler

	findExamInfos gt.Handler

	pb.UnimplementedExamServiceServer
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpointds endpoint.Endpoints, logger log.Logger) pb.ExamServiceServer {
	return &GRPCServer{
		logger: logger,

		// Exam
		createExam: gt.NewServer(
			endpointds.CreateExam,
			decodeCreateExamRequest,
			encodeCreateExamResponse,
		),
		updateExam: gt.NewServer(
			endpointds.UpdateExam,
			decodeUpdateExamRequest,
			encodeUpdateExamResponse,
		),
		findExams: gt.NewServer(
			endpointds.FindExams,
			decodeFindExamsRequest,
			encodeFindExamsResponse,
		),
		deleteExam: gt.NewServer(
			endpointds.DeleteExam,
			decodeDeleteExamRequest,
			encodeDeleteExamResponse,
		),

		// Question
		createQuestion: gt.NewServer(
			endpointds.CreateQuestion,
			decodeCreateQuestionRequest,
			encodeCreateQuestionResponse,
		),
		updateQuestion: gt.NewServer(
			endpointds.UpdateQuestion,
			decodeUpdateQuestionRequest,
			encodeUpdateQuestionResponse,
		),
		findQuestions: gt.NewServer(
			endpointds.FindQuestions,
			decodeFindQuestionsRequest,
			encodeFindQuestionsResponse,
		),
		deleteQuestion: gt.NewServer(
			endpointds.DeleteQuestion,
			decodeDeleteQuestionRequest,
			encodeDeleteQuestionResponse,
		),
		findRandomQuestions: gt.NewServer(
			endpointds.FindRandomQuestions,
			decodeFindRandomQuestionsRequest,
			encodeFindRandomQuestionsResponse,
		),

		// ExamRecord
		createExamRecord: gt.NewServer(
			endpointds.CreateExamRecord,
			decodeCreateExamRecordRequest,
			encodeCreateExamRecordResponse,
		),
		findExamRecords: gt.NewServer(
			endpointds.FindExamRecords,
			decodeFindExamRecordsRequest,
			encodeFindExamRecordsResponse,
		),
		findExamRecordOverview: gt.NewServer(
			endpointds.FindExamRecordOverview,
			decodeFindExamRecordOverviewRequest,
			encodeFindExamRecordOverviewResponse,
		),

		// ExamInfo
		findExamInfos: gt.NewServer(
			endpointds.FindExamInfos,
			decodeFindExamInfosRequest,
			encodeFindExamInfosResponse,
		),
	}
}

func (s GRPCServer) CreateExam(
	ctx context.Context,
	req *pb.CreateExamRequest,
) (*pb.CreateExamResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.createExam.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CreateExamResponse), nil
}

func decodeCreateExamRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.CreateExamRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.CreateExamRequest{
		Topic:       req.Topic,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		UserId:      req.UserId,
	}, nil
}

func encodeCreateExamResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.CreateExamResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.CreateExamResponse{
		ExamId: resp.ExamId,
	}, nil
}

func (s GRPCServer) UpdateExam(
	ctx context.Context,
	req *pb.UpdateExamRequest,
) (*pb.UpdateExamResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.updateExam.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.UpdateExamResponse), nil
}

func decodeUpdateExamRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.UpdateExamRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.UpdateExamRequest{
		ExamId:      req.ExamId,
		Topic:       req.Topic,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		UserId:      req.UserId,
	}, nil
}

func encodeUpdateExamResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.UpdateExamResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.UpdateExamResponse{
		ExamId: resp.ExamId,
	}, nil
}

func (s GRPCServer) FindExams(
	ctx context.Context,
	req *pb.FindExamsRequest,
) (*pb.FindExamsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.findExams.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindExamsResponse), nil
}

func decodeFindExamsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.FindExamsRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindExamsRequest{
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
		UserId:    req.UserId,
	}, nil
}

func encodeFindExamsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.FindExamsResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	exams := []*pb.Exam{}

	for _, exam := range resp.Exams {
		exams = append(exams, toPBExam(&exam))
	}

	return &pb.FindExamsResponse{
		Total:     resp.Total,
		PageCount: resp.PageCount,
		Exams:     exams,
	}, nil
}

func (s GRPCServer) DeleteExam(
	ctx context.Context,
	req *pb.DeleteExamRequest,
) (*pb.DeleteExamResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.deleteExam.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.DeleteExamResponse), nil
}

func decodeDeleteExamRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.DeleteExamRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.DeleteExamRequest{
		ExamId: req.ExamId,
		UserId: req.UserId,
	}, nil
}

func encodeDeleteExamResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(endpoint.DeleteExamResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.DeleteExamResponse{}, nil
}

func (s GRPCServer) CreateQuestion(
	ctx context.Context,
	req *pb.CreateQuestionRequest,
) (*pb.CreateQuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.createQuestion.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CreateQuestionResponse), nil
}

func decodeCreateQuestionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.CreateQuestionRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.CreateQuestionRequest{
		ExamId:  req.ExamId,
		Ask:     req.Ask,
		Answers: req.Answers,
		UserId:  req.UserId,
	}, nil
}

func encodeCreateQuestionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.CreateQuestionResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.CreateQuestionResponse{
		QuestionId: resp.QuestionId,
	}, nil
}

func (s GRPCServer) UpdateQuestion(
	ctx context.Context,
	req *pb.UpdateQuestionRequest,
) (*pb.UpdateQuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.updateQuestion.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.UpdateQuestionResponse), nil
}

func decodeUpdateQuestionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.UpdateQuestionRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.UpdateQuestionRequest{
		QuestionId: req.QuestionId,
		Ask:        req.Ask,
		Answers:    req.Answers,
		UserId:     req.UserId,
	}, nil
}

func encodeUpdateQuestionResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.UpdateQuestionResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.UpdateQuestionResponse{
		QuestionId: resp.QuestionId,
	}, nil
}

func (s GRPCServer) FindQuestions(
	ctx context.Context,
	req *pb.FindQuestionsRequest,
) (*pb.FindQuestionsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.findQuestions.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindQuestionsResponse), nil
}

func decodeFindQuestionsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.FindQuestionsRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindQuestionsRequest{
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
		ExamId:    req.ExamId,
		UserId:    req.UserId,
	}, nil
}

func encodeFindQuestionsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.FindQuestionsResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	questions := []*pb.Question{}

	for _, question := range resp.Questions {
		questions = append(questions, toPBQuestion(&question))
	}

	return &pb.FindQuestionsResponse{
		Total:     resp.Total,
		PageCount: resp.PageCount,
		Questions: questions,
	}, nil
}

func (s GRPCServer) DeleteQuestion(
	ctx context.Context,
	req *pb.DeleteQuestionRequest,
) (*pb.DeleteQuestionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.deleteQuestion.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.DeleteQuestionResponse), nil
}

func decodeDeleteQuestionRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.DeleteQuestionRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.DeleteQuestionRequest{
		QuestionId: req.QuestionId,
		UserId:     req.UserId,
	}, nil
}

func encodeDeleteQuestionResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(endpoint.DeleteQuestionResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.DeleteQuestionResponse{}, nil
}

func (s GRPCServer) FindRandomQuestions(
	ctx context.Context,
	req *pb.FindRandomQuestionsRequest,
) (*pb.FindRandomQuestionsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.findRandomQuestions.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindRandomQuestionsResponse), nil
}

func decodeFindRandomQuestionsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.FindRandomQuestionsRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindRandomQuestionsRequest{
		ExamId: req.ExamId,
		UserId: req.UserId,
		Size:   req.Size,
	}, nil
}

func encodeFindRandomQuestionsResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	resp, ok := response.(endpoint.FindRandomQuestionsResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	exam := toPBExam(resp.Exam)
	questions := []*pb.Question{}

	for _, question := range resp.Questions {
		questions = append(questions, toPBQuestion(&question))
	}

	return &pb.FindRandomQuestionsResponse{
		Exam:      exam,
		Questions: questions,
	}, nil
}

func (s GRPCServer) CreateExamRecord(
	ctx context.Context,
	req *pb.CreateExamRecordRequest,
) (*pb.CreateExamRecordResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.createExamRecord.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CreateExamRecordResponse), nil
}

func decodeCreateExamRecordRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.CreateExamRecordRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.CreateExamRecordRequest{
		ExamId:           req.ExamId,
		Score:            req.Score,
		WrongQuestionIds: req.WrongQuestionIds,
		UserId:           req.UserId,
	}, nil
}

func encodeCreateExamRecordResponse(_ context.Context, response interface{}) (interface{}, error) {
	_, ok := response.(endpoint.CreateExamRecordResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.CreateExamRecordResponse{}, nil
}

func (s GRPCServer) FindExamRecords(
	ctx context.Context,
	req *pb.FindExamRecordsRequest,
) (*pb.FindExamRecordsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.findExamRecords.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindExamRecordsResponse), nil
}

func decodeFindExamRecordsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.FindExamRecordsRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindExamRecordsRequest{
		PageIndex: req.PageIndex,
		PageSize:  req.PageSize,
		ExamId:    req.ExamId,
		UserId:    req.UserId,
	}, nil
}

func encodeFindExamRecordsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.FindExamRecordsResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	examRecords := []*pb.ExamRecord{}

	for _, examRecord := range resp.ExamRecords {
		examRecords = append(examRecords, toPBExamRecord(&examRecord))
	}

	return &pb.FindExamRecordsResponse{
		Total:       resp.Total,
		PageCount:   resp.PageCount,
		ExamRecords: examRecords,
	}, nil
}

func (s GRPCServer) FindExamRecordOverview(
	ctx context.Context,
	req *pb.FindExamRecordOverviewRequest,
) (*pb.FindExamRecordOverviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.findExamRecordOverview.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindExamRecordOverviewResponse), nil
}

func decodeFindExamRecordOverviewRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req, ok := request.(*pb.FindExamRecordOverviewRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindExamRecordOverviewRequest{
		ExamId:    req.ExamId,
		UserId:    req.UserId,
		StartDate: req.StartDate.AsTime(),
	}, nil
}

func encodeFindExamRecordOverviewResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	resp, ok := response.(endpoint.FindExamRecordOverviewResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	exam := toPBExam(resp.Exam)
	questions := []*pb.Question{}

	for _, question := range resp.Questions {
		questions = append(questions, toPBQuestion(&question))
	}

	answerWrongs := []*pb.AnswerWrong{}

	for _, answerWrong := range resp.AnswerWrongs {
		answerWrongs = append(answerWrongs, toPBAnswerWrong(&answerWrong))
	}

	examRecords := []*pb.ExamRecord{}

	for _, examRecord := range resp.ExamRecords {
		examRecords = append(examRecords, toPBExamRecord(&examRecord))
	}

	return &pb.FindExamRecordOverviewResponse{
		StartDate:    resp.StartDate,
		Exam:         exam,
		Questions:    questions,
		AnswerWrongs: answerWrongs,
		ExamRecords:  examRecords,
	}, nil
}

func (s GRPCServer) FindExamInfos(
	ctx context.Context,
	req *pb.FindExamInfosRequest,
) (*pb.FindExamInfosResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()
	_, resp, err := s.findExamInfos.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindExamInfosResponse), nil
}

func decodeFindExamInfosRequest(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*pb.FindExamInfosRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindExamInfosRequest{
		UserId:   req.UserId,
		IsPublic: req.IsPublic,
	}, nil
}

func encodeFindExamInfosResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(endpoint.FindExamInfosResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	examInfos := []*pb.ExamInfo{}

	for _, examInfo := range resp.ExamInfos {
		examInfos = append(examInfos, &pb.ExamInfo{
			ExamId:        examInfo.ExamId,
			Topic:         examInfo.Topic,
			Description:   examInfo.Description,
			IsPublic:      examInfo.IsPublic,
			QuestionCount: examInfo.QuestionCount,
			RecordCount:   examInfo.RecordCount,
		})
	}

	return &pb.FindExamInfosResponse{
		ExamInfos: examInfos,
	}, nil
}

func toPBExam(exam *model.Exam) *pb.Exam {
	if exam == nil {
		return nil
	}

	return &pb.Exam{
		Id:          exam.Id.Hex(),
		Topic:       exam.Topic,
		Description: exam.Description,
		Tags:        exam.Tags,
		IsPublic:    exam.IsPublic,
		UserId:      exam.UserId,
		CreatedAt:   timestamppb.New(exam.CreatedAt),
		UpdatedAt:   timestamppb.New(exam.UpdatedAt),
	}
}

func toPBQuestion(question *model.Question) *pb.Question {
	if question == nil {
		return nil
	}

	return &pb.Question{
		Id:        question.Id.Hex(),
		ExamId:    question.ExamId,
		Ask:       question.Ask,
		Answers:   question.Answers,
		UserId:    question.UserId,
		CreatedAt: timestamppb.New(question.CreatedAt),
		UpdatedAt: timestamppb.New(question.UpdatedAt),
	}
}

func toPBAnswerWrong(answerWrong *model.AnswerWrong) *pb.AnswerWrong {
	if answerWrong == nil {
		return nil
	}

	return &pb.AnswerWrong{
		Id:         answerWrong.Id.Hex(),
		ExamId:     answerWrong.ExamId,
		QuestionId: answerWrong.QuestionId,
		Times:      answerWrong.Times,
		UserId:     answerWrong.UserId,
		CreatedAt:  timestamppb.New(answerWrong.CreatedAt),
		UpdatedAt:  timestamppb.New(answerWrong.UpdatedAt),
	}
}

func toPBExamRecord(examRecord *model.ExamRecord) *pb.ExamRecord {
	if examRecord == nil {
		return nil
	}

	return &pb.ExamRecord{
		Id:        examRecord.Id.Hex(),
		ExamId:    examRecord.ExamId,
		Score:     examRecord.Score,
		UserId:    examRecord.UserId,
		CreatedAt: timestamppb.New(examRecord.CreatedAt),
		UpdatedAt: timestamppb.New(examRecord.UpdatedAt),
	}
}
