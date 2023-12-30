package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/service"
)

type Endpoints struct {
	FindWordByDictionary      endpoint.Endpoint
	CreateFavoriteWordMeaning endpoint.Endpoint
	DeleteFavoriteWordMeaning endpoint.Endpoint
	FindFavoriteWordMeanings  endpoint.Endpoint
}

// MakeAddEndpoint struct holds the endpoint response definition
func MakeEndpoints(wordService service.WordService, logger log.Logger) Endpoints {
	findWordByDictionaryEndpoint := makeFindWordByDictionaryEndpoint(wordService)
	findWordByDictionaryEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindWordByDictionary"))(findWordByDictionaryEndpoint)

	createFavoriteWordMeaningEndpoint := makeCreateFavoriteWordMeaningEndpoint(wordService)
	createFavoriteWordMeaningEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateFavoriteWordMeaning"))(createFavoriteWordMeaningEndpoint)

	deleteFavoriteWordMeaningEndpoint := makeDeleteFavoriteWordMeaningEndpoint(wordService)
	deleteFavoriteWordMeaningEndpoint = LoggingMiddleware(
		log.With(logger, "method", "DeleteFavoriteWordMeaning"))(deleteFavoriteWordMeaningEndpoint)

	findFavoriteWordMeaningsEndpoint := makeFindFavoriteWordMeaningsEndpoint(wordService)
	findFavoriteWordMeaningsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindFavoriteWordMeanings"))(findFavoriteWordMeaningsEndpoint)

	return Endpoints{
		FindWordByDictionary:      findWordByDictionaryEndpoint,
		CreateFavoriteWordMeaning: createFavoriteWordMeaningEndpoint,
		DeleteFavoriteWordMeaning: deleteFavoriteWordMeaningEndpoint,
		FindFavoriteWordMeanings:  findFavoriteWordMeaningsEndpoint,
	}
}

type FindWordByDictionaryRequest struct {
	Word   string
	UserId string
}

type FindWordByDictionaryResponse struct {
	WordMeanings []model.WordMeaning
}

func makeFindWordByDictionaryEndpoint(wordService service.WordService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindWordByDictionaryRequest)
		wordMeangins, err := wordService.FindWordByDictionary(req.Word, req.UserId)
		if err != nil {
			return nil, err
		}
		return FindWordByDictionaryResponse{WordMeanings: wordMeangins}, nil
	}
}

type CreateFavoriteWordMeaningRequest struct {
	UserId        string
	WordMeaningId string
}

type CreateFavoriteWordMeaningResponse struct {
	FavoriteWordMeaningId string
}

func makeCreateFavoriteWordMeaningEndpoint(wordService service.WordService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateFavoriteWordMeaningRequest)
		favoriteWordMeaningId, err := wordService.CreateFavoriteWordMeaning(
			req.UserId,
			req.WordMeaningId,
		)
		if err != nil {
			return nil, err
		}
		return CreateFavoriteWordMeaningResponse{
			FavoriteWordMeaningId: favoriteWordMeaningId,
		}, nil
	}
}

type DeleteFavoriteWordMeaningRequest struct {
	FavoriteWordMeaningId string
	UserId                string
}

type DeleteFavoriteWordMeaningResponse struct{}

func makeDeleteFavoriteWordMeaningEndpoint(wordService service.WordService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteFavoriteWordMeaningRequest)
		err := wordService.DeleteFavoriteWordMeaning(
			req.FavoriteWordMeaningId,
			req.UserId,
		)
		if err != nil {
			return nil, err
		}
		return DeleteFavoriteWordMeaningResponse{}, nil
	}
}

type FindFavoriteWordMeaningsRequest struct {
	PageInde int64
	PageSize int64
	UserId   string
	Word     string
}

type FindFavoriteWordMeaningsResponse struct {
	Total                int64
	PageCount            int64
	FavoriteWordMeanings []model.WordMeaning
}

func makeFindFavoriteWordMeaningsEndpoint(wordService service.WordService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindFavoriteWordMeaningsRequest)
		total, pageCount, favoriteWordMeanings, err := wordService.FindFavoriteWordMeanings(
			req.PageInde,
			req.PageSize,
			req.UserId,
			req.Word,
		)
		if err != nil {
			return nil, err
		}
		return FindFavoriteWordMeaningsResponse{
			Total:                total,
			PageCount:            pageCount,
			FavoriteWordMeanings: favoriteWordMeanings,
		}, nil
	}
}
