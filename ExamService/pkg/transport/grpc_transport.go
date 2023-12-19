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
	pb.UnimplementedExamServiceServer
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpointds endpoint.Endpoints, logger log.Logger) pb.ExamServiceServer {
	return &GRPCServer{
		createExam: gt.NewServer(
			endpointds.CreateExam,
			decodeCreateExamRequest,
			encodeCreateExamResponse,
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
