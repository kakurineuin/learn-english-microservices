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
}

// MakeAddEndpoint struct holds the endpoint response definition
func MakeEndpoints(wordService service.WordService, logger log.Logger) Endpoints {
	findWordByDictionaryEndpoint := makeFindWordByDictionaryEndpoint(wordService)
	findWordByDictionaryEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindWordByDictionary"))(findWordByDictionaryEndpoint)

	createFavoriteWordMeaningEndpoint := makeCreateFavoriteWordMeaningEndpoint(wordService)
	createFavoriteWordMeaningEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateFavoriteWordMeaning"))(createFavoriteWordMeaningEndpoint)

	return Endpoints{
		FindWordByDictionary:      findWordByDictionaryEndpoint,
		CreateFavoriteWordMeaning: createFavoriteWordMeaningEndpoint,
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
