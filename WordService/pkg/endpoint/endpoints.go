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
	var findWordByDictionaryEndpoint endpoint.Endpoint
	{
		findWordByDictionaryEndpoint = makeFindWordByDictionaryEndpoint(wordService)
		findWordByDictionaryEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findWordByDictionaryEndpoint,
		)
		findWordByDictionaryEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findWordByDictionaryEndpoint,
		)
		findWordByDictionaryEndpoint = LoggingMiddleware(
			log.With(logger, "method", "FindWordByDictionary"))(findWordByDictionaryEndpoint)
		findWordByDictionaryEndpoint = RecoverMiddleware(
			log.With(logger, "method", "FindWordByDictionary"))(findWordByDictionaryEndpoint)
	}

	var createFavoriteWordMeaningEndpoint endpoint.Endpoint
	{
		createFavoriteWordMeaningEndpoint = makeCreateFavoriteWordMeaningEndpoint(wordService)
		createFavoriteWordMeaningEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			createFavoriteWordMeaningEndpoint,
		)
		createFavoriteWordMeaningEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			createFavoriteWordMeaningEndpoint,
		)
		createFavoriteWordMeaningEndpoint = LoggingMiddleware(
			log.With(
				logger,
				"method",
				"CreateFavoriteWordMeaning",
			),
		)(
			createFavoriteWordMeaningEndpoint,
		)
		createFavoriteWordMeaningEndpoint = RecoverMiddleware(
			log.With(
				logger,
				"method",
				"CreateFavoriteWordMeaning",
			),
		)(
			createFavoriteWordMeaningEndpoint,
		)
	}

	var deleteFavoriteWordMeaningEndpoint endpoint.Endpoint
	{
		deleteFavoriteWordMeaningEndpoint = makeDeleteFavoriteWordMeaningEndpoint(wordService)
		deleteFavoriteWordMeaningEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			deleteFavoriteWordMeaningEndpoint,
		)
		deleteFavoriteWordMeaningEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			deleteFavoriteWordMeaningEndpoint,
		)
		deleteFavoriteWordMeaningEndpoint = LoggingMiddleware(
			log.With(
				logger,
				"method",
				"DeleteFavoriteWordMeaning",
			),
		)(
			deleteFavoriteWordMeaningEndpoint,
		)
		deleteFavoriteWordMeaningEndpoint = RecoverMiddleware(
			log.With(
				logger,
				"method",
				"DeleteFavoriteWordMeaning",
			),
		)(
			deleteFavoriteWordMeaningEndpoint,
		)
	}

	var findFavoriteWordMeaningsEndpoint endpoint.Endpoint
	{
		findFavoriteWordMeaningsEndpoint = makeFindFavoriteWordMeaningsEndpoint(wordService)
		findFavoriteWordMeaningsEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findFavoriteWordMeaningsEndpoint,
		)
		findFavoriteWordMeaningsEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findFavoriteWordMeaningsEndpoint,
		)
		findFavoriteWordMeaningsEndpoint = LoggingMiddleware(
			log.With(
				logger,
				"method",
				"FindFavoriteWordMeanings",
			),
		)(
			findFavoriteWordMeaningsEndpoint,
		)
		findFavoriteWordMeaningsEndpoint = RecoverMiddleware(
			log.With(
				logger,
				"method",
				"FindFavoriteWordMeanings",
			),
		)(
			findFavoriteWordMeaningsEndpoint,
		)
	}

	var findRandomFavoriteWordMeaningsEndpoint endpoint.Endpoint
	{
		findRandomFavoriteWordMeaningsEndpoint = makeFindRandomFavoriteWordMeaningsEndpoint(
			wordService,
		)
		findRandomFavoriteWordMeaningsEndpoint = ratelimit.NewErroringLimiter(
			rate.NewLimiter(rate.Every(time.Second), 1),
		)(
			findRandomFavoriteWordMeaningsEndpoint,
		)
		findRandomFavoriteWordMeaningsEndpoint = circuitbreaker.Gobreaker(
			gobreaker.NewCircuitBreaker(gobreaker.Settings{}),
		)(
			findRandomFavoriteWordMeaningsEndpoint,
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
	}

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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindWordByDictionaryRequest)
		wordMeangins, err := wordService.FindWordByDictionary(ctx, req.Word, req.UserId)
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateFavoriteWordMeaningRequest)
		favoriteWordMeaningId, err := wordService.CreateFavoriteWordMeaning(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteFavoriteWordMeaningRequest)
		err := wordService.DeleteFavoriteWordMeaning(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindFavoriteWordMeaningsRequest)
		total, pageCount, favoriteWordMeanings, err := wordService.FindFavoriteWordMeanings(
			ctx,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRandomFavoriteWordMeaningsRequest)
		favoriteWordMeanings, err := wordService.FindRandomFavoriteWordMeanings(
			ctx,
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
