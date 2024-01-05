package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/service"
)

type Endpoints struct {
	FindWordByDictionary           endpoint.Endpoint
	CreateFavoriteWordMeaning      endpoint.Endpoint
	DeleteFavoriteWordMeaning      endpoint.Endpoint
	FindFavoriteWordMeanings       endpoint.Endpoint
	FindRandomFavoriteWordMeanings endpoint.Endpoint
}

// MakeAddEndpoint struct holds the endpoint response definition
func MakeEndpoints(wordService service.WordService, logger log.Logger) Endpoints {
	findWordByDictionaryEndpoint := makeFindWordByDictionaryEndpoint(wordService)
	findWordByDictionaryEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindWordByDictionary"))(findWordByDictionaryEndpoint)
	findWordByDictionaryEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindWordByDictionary"))(findWordByDictionaryEndpoint)

	createFavoriteWordMeaningEndpoint := makeCreateFavoriteWordMeaningEndpoint(wordService)
	createFavoriteWordMeaningEndpoint = LoggingMiddleware(
		log.With(logger, "method", "CreateFavoriteWordMeaning"))(createFavoriteWordMeaningEndpoint)
	createFavoriteWordMeaningEndpoint = RecoverMiddleware(
		log.With(logger, "method", "CreateFavoriteWordMeaning"))(createFavoriteWordMeaningEndpoint)

	deleteFavoriteWordMeaningEndpoint := makeDeleteFavoriteWordMeaningEndpoint(wordService)
	deleteFavoriteWordMeaningEndpoint = LoggingMiddleware(
		log.With(logger, "method", "DeleteFavoriteWordMeaning"))(deleteFavoriteWordMeaningEndpoint)
	deleteFavoriteWordMeaningEndpoint = RecoverMiddleware(
		log.With(logger, "method", "DeleteFavoriteWordMeaning"))(deleteFavoriteWordMeaningEndpoint)

	findFavoriteWordMeaningsEndpoint := makeFindFavoriteWordMeaningsEndpoint(wordService)
	findFavoriteWordMeaningsEndpoint = LoggingMiddleware(
		log.With(logger, "method", "FindFavoriteWordMeanings"))(findFavoriteWordMeaningsEndpoint)
	findFavoriteWordMeaningsEndpoint = RecoverMiddleware(
		log.With(logger, "method", "FindFavoriteWordMeanings"))(findFavoriteWordMeaningsEndpoint)

	findRandomFavoriteWordMeaningsEndpoint := makeFindRandomFavoriteWordMeaningsEndpoint(
		wordService,
	)
	findRandomFavoriteWordMeaningsEndpoint = LoggingMiddleware(
		log.With(
			logger,
			"method",
			"FindRandomFavoriteWordMeanings",
		),
	)(
		findRandomFavoriteWordMeaningsEndpoint,
	)
	findRandomFavoriteWordMeaningsEndpoint = RecoverMiddleware(
		log.With(
			logger,
			"method",
			"FindRandomFavoriteWordMeanings",
		),
	)(
		findRandomFavoriteWordMeaningsEndpoint,
	)

	return Endpoints{
		FindWordByDictionary:           findWordByDictionaryEndpoint,
		CreateFavoriteWordMeaning:      createFavoriteWordMeaningEndpoint,
		DeleteFavoriteWordMeaning:      deleteFavoriteWordMeaningEndpoint,
		FindFavoriteWordMeanings:       findFavoriteWordMeaningsEndpoint,
		FindRandomFavoriteWordMeanings: findRandomFavoriteWordMeaningsEndpoint,
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
	PageInde int32
	PageSize int32
	UserId   string
	Word     string
}

type FindFavoriteWordMeaningsResponse struct {
	Total                int32
	PageCount            int32
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

type FindRandomFavoriteWordMeaningsRequest struct {
	UserId string
	Size   int32
}

type FindRandomFavoriteWordMeaningsResponse struct {
	FavoriteWordMeanings []model.WordMeaning
}

func makeFindRandomFavoriteWordMeaningsEndpoint(wordService service.WordService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRandomFavoriteWordMeaningsRequest)
		favoriteWordMeanings, err := wordService.FindRandomFavoriteWordMeanings(
			req.UserId,
			req.Size,
		)
		if err != nil {
			return nil, err
		}
		return FindRandomFavoriteWordMeaningsResponse{
			FavoriteWordMeanings: favoriteWordMeanings,
		}, nil
	}
}
