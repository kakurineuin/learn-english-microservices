package transport

import (
	"context"
	"errors"

	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/word-service/pb"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/endpoint"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

type GRPCServer struct {
	findWordByDictionary      gt.Handler
	createFavoriteWordMeaning gt.Handler
	deleteFavoriteWordMeaning gt.Handler
	findFavoriteWordMeanings  gt.Handler

	pb.UnimplementedWordServiceServer
}

// NewGRPCServer initializes a new gRPC server
func NewGRPCServer(endpointds endpoint.Endpoints, logger log.Logger) pb.WordServiceServer {
	return &GRPCServer{
		findWordByDictionary: gt.NewServer(
			endpointds.FindWordByDictionary,
			decodeFindWordByDictionaryRequest,
			encodeFindWordByDictionaryResponse,
		),
		createFavoriteWordMeaning: gt.NewServer(
			endpointds.CreateFavoriteWordMeaning,
			decodeCreateFavoriteWordMeaningRequest,
			encodeCreateFavoriteWordMeaningResponse,
		),
		deleteFavoriteWordMeaning: gt.NewServer(
			endpointds.DeleteFavoriteWordMeaning,
			decodeDeleteFavoriteWordMeaningRequest,
			encodeDeleteFavoriteWordMeaningResponse,
		),
		findFavoriteWordMeanings: gt.NewServer(
			endpointds.FindFavoriteWordMeanings,
			decodeFindFavoriteWordMeaningsRequest,
			encodeFindFavoriteWordMeaningsResponse,
		),
	}
}

func (s GRPCServer) FindWordByDictionary(
	ctx context.Context,
	req *pb.FindWordByDictionaryRequest,
) (*pb.FindWordByDictionaryResponse, error) {
	_, resp, err := s.findWordByDictionary.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindWordByDictionaryResponse), nil
}

func decodeFindWordByDictionaryRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req, ok := request.(*pb.FindWordByDictionaryRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindWordByDictionaryRequest{
		Word:   req.Word,
		UserId: req.UserId,
	}, nil
}

func encodeFindWordByDictionaryResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	resp, ok := response.(endpoint.FindWordByDictionaryResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.FindWordByDictionaryResponse{
		WordMeanings: toPBWordMeanings(resp.WordMeanings),
	}, nil
}

func (s GRPCServer) CreateFavoriteWordMeaning(
	ctx context.Context,
	req *pb.CreateFavoriteWordMeaningRequest,
) (*pb.CreateFavoriteWordMeaningResponse, error) {
	_, resp, err := s.createFavoriteWordMeaning.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.CreateFavoriteWordMeaningResponse), nil
}

func decodeCreateFavoriteWordMeaningRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req, ok := request.(*pb.CreateFavoriteWordMeaningRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.CreateFavoriteWordMeaningRequest{
		UserId:        req.UserId,
		WordMeaningId: req.WordMeaningId,
	}, nil
}

func encodeCreateFavoriteWordMeaningResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	resp, ok := response.(endpoint.CreateFavoriteWordMeaningResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.CreateFavoriteWordMeaningResponse{
		FavoriteWordMeaningId: resp.FavoriteWordMeaningId,
	}, nil
}

func (s GRPCServer) DeleteFavoriteWordMeaning(
	ctx context.Context,
	req *pb.DeleteFavoriteWordMeaningRequest,
) (*pb.DeleteFavoriteWordMeaningResponse, error) {
	_, resp, err := s.deleteFavoriteWordMeaning.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.DeleteFavoriteWordMeaningResponse), nil
}

func decodeDeleteFavoriteWordMeaningRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req, ok := request.(*pb.DeleteFavoriteWordMeaningRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.DeleteFavoriteWordMeaningRequest{
		FavoriteWordMeaningId: req.FavoriteWordMeaningId,
		UserId:                req.UserId,
	}, nil
}

func encodeDeleteFavoriteWordMeaningResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	_, ok := response.(endpoint.DeleteFavoriteWordMeaningResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.DeleteFavoriteWordMeaningResponse{}, nil
}

func (s GRPCServer) FindFavoriteWordMeanings(
	ctx context.Context,
	req *pb.FindFavoriteWordMeaningsRequest,
) (*pb.FindFavoriteWordMeaningsResponse, error) {
	_, resp, err := s.findFavoriteWordMeanings.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.FindFavoriteWordMeaningsResponse), nil
}

func decodeFindFavoriteWordMeaningsRequest(
	_ context.Context,
	request interface{},
) (interface{}, error) {
	req, ok := request.(*pb.FindFavoriteWordMeaningsRequest)
	if !ok {
		return nil, errors.New("invalid request body")
	}

	return endpoint.FindFavoriteWordMeaningsRequest{
		PageInde: req.PageIndex,
		PageSize: req.PageSize,
		UserId:   req.UserId,
		Word:     req.Word,
	}, nil
}

func encodeFindFavoriteWordMeaningsResponse(
	_ context.Context,
	response interface{},
) (interface{}, error) {
	resp, ok := response.(endpoint.FindFavoriteWordMeaningsResponse)
	if !ok {
		return nil, errors.New("invalid response body")
	}

	return &pb.FindFavoriteWordMeaningsResponse{
		Total:                resp.Total,
		PageCount:            resp.PageCount,
		FavoriteWordMeanings: toPBWordMeanings(resp.FavoriteWordMeanings),
	}, nil
}

func toPBWordMeanings(wordMeanings []model.WordMeaning) []*pb.WordMeaning {
	pbWordMeanings := []*pb.WordMeaning{}

	for _, wm := range wordMeanings {
		pbExamples := []*pb.Example{}

		for _, example := range wm.Examples {
			pbSentences := []*pb.Sentence{}

			for _, sentence := range example.Examples {
				pbSentences = append(pbSentences, &pb.Sentence{
					AudioUrl: sentence.AudioUrl,
					Text:     sentence.Text,
				})
			}

			pbExamples = append(pbExamples, &pb.Example{
				Pattern:  example.Pattern,
				Examples: pbSentences,
			})
		}

		id := ""

		if !wm.Id.IsZero() {
			id = wm.Id.Hex()
		}

		favoriteWordMeaningId := ""

		if !wm.FavoriteWordMeaningId.IsZero() {
			favoriteWordMeaningId = wm.FavoriteWordMeaningId.Hex()
		}

		pbWordMeaning := pb.WordMeaning{
			Id:           id,
			Word:         wm.Word,
			PartOfSpeech: wm.PartOfSpeech,
			Gram:         wm.Gram,
			Pronunciation: &pb.Pronunciation{
				Text:       wm.Pronunciation.Text,
				UkAudioUrl: wm.Pronunciation.UkAudioUrl,
				UsAudioUrl: wm.Pronunciation.UsAudioUrl,
			},
			DefGram:               wm.DefGram,
			Definition:            wm.Definition,
			Examples:              pbExamples,
			OrderByNo:             wm.OrderByNo,
			QueryByWords:          wm.QueryByWords,
			FavoriteWordMeaningId: favoriteWordMeaningId,
		}
		pbWordMeanings = append(pbWordMeanings, &pbWordMeaning)
	}

	return pbWordMeanings
}
