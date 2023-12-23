package transport

import (
	"context"
	"errors"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pb"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/endpoint"
)

type GRPCServer struct {
	createExam gt.Handler
	updateExam gt.Handler
	findExams  gt.Handler
	deleteExam gt.Handler

	createQuestion gt.Handler
	updateQuestion gt.Handler
	findQuestions  gt.Handler
	deleteQuestion gt.Handler

	createExamRecord gt.Handler
	findExamRecords  gt.Handler

	pb.UnimplementedExamServiceServer
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpointds endpoint.Endpoints, logger log.Logger) pb.ExamServiceServer {
	return &GRPCServer{
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
	}
}

func (s GRPCServer) CreateExam(
	ctx context.Context,
	req *pb.CreateExamRequest,
) (*pb.CreateExamResponse, error) {
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
		exams = append(exams, &pb.Exam{
			Id:          exam.Id.Hex(),
			Topic:       exam.Topic,
			Description: exam.Description,
			IsPublic:    exam.IsPublic,
			Tags:        exam.Tags,
			UserId:      exam.UserId,
		})
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
		questions = append(questions, &pb.Question{
			Id:      question.Id.Hex(),
			ExamId:  question.ExamId,
			Ask:     question.Ask,
			Answers: question.Answers,
			UserId:  question.UserId,
		})
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

func (s GRPCServer) CreateExamRecord(
	ctx context.Context,
	req *pb.CreateExamRecordRequest,
) (*pb.CreateExamRecordResponse, error) {
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
		examRecords = append(examRecords, &pb.ExamRecord{
			Id:     examRecord.Id.Hex(),
			ExamId: examRecord.ExamId,
			Score:  examRecord.Score,
			UserId: examRecord.UserId,
		})
	}

	return &pb.FindExamRecordsResponse{
		Total:       resp.Total,
		PageCount:   resp.PageCount,
		ExamRecords: examRecords,
	}, nil
}
