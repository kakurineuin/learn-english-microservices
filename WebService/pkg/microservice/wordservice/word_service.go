package wordservice

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
)

const SERVER_ADDRESS = "localhost:8091"

//go:generate mockery --name WordService
type WordService interface {
	Connect() error
	Disconnect() error
	FindWordByDictionary(
		word, userId string,
	) (*pb.FindWordByDictionaryResponse, error)
	CreateFavoriteWordMeaning(
		userId, wordMeaningId string,
	) (*pb.CreateFavoriteWordMeaningResponse, error)
	DeleteFavoriteWordMeaning(
		favoriteWordMeaningId, userId string,
	) (*pb.DeleteFavoriteWordMeaningResponse, error)
	FindFavoriteWordMeanings(
		pageIndex, pageSize int32,
		userId, word string,
	) (*pb.FindFavoriteWordMeaningsResponse, error)
	FindRandomFavoriteWordMeanings(
		userId string, size int32,
	) (*pb.FindRandomFavoriteWordMeaningsResponse, error)
}

func New() WordService {
	return &wordService{}
}

type wordService struct {
	connection *grpc.ClientConn
	client     pb.WordServiceClient
}

func (service *wordService) Connect() error {
	conn, err := grpc.Dial(
		SERVER_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	service.connection = conn
	service.client = pb.NewWordServiceClient(conn)

	return nil
}

func (service *wordService) Disconnect() error {
	if err := service.connection.Close(); err != nil {
		return err
	}

	service.connection = nil
	service.client = nil

	return nil
}

func (service wordService) FindWordByDictionary(
	word, userId string,
) (*pb.FindWordByDictionaryResponse, error) {
	return service.client.FindWordByDictionary(
		context.Background(),
		&pb.FindWordByDictionaryRequest{
			Word:   word,
			UserId: userId,
		},
	)
}

func (service wordService) CreateFavoriteWordMeaning(
	userId, wordMeaningId string,
) (*pb.CreateFavoriteWordMeaningResponse, error) {
	return service.client.CreateFavoriteWordMeaning(
		context.Background(),
		&pb.CreateFavoriteWordMeaningRequest{
			UserId:        userId,
			WordMeaningId: wordMeaningId,
		},
	)
}

func (service wordService) DeleteFavoriteWordMeaning(
	favoriteWordMeaningId, userId string,
) (*pb.DeleteFavoriteWordMeaningResponse, error) {
	return service.client.DeleteFavoriteWordMeaning(
		context.Background(),
		&pb.DeleteFavoriteWordMeaningRequest{
			FavoriteWordMeaningId: favoriteWordMeaningId,
			UserId:                userId,
		},
	)
}

func (service wordService) FindFavoriteWordMeanings(
	pageIndex, pageSize int32,
	userId, word string,
) (*pb.FindFavoriteWordMeaningsResponse, error) {
	return service.client.FindFavoriteWordMeanings(
		context.Background(),
		&pb.FindFavoriteWordMeaningsRequest{
			PageIndex: pageIndex,
			PageSize:  pageSize,
			UserId:    userId,
			Word:      word,
		},
	)
}

func (service wordService) FindRandomFavoriteWordMeanings(
	userId string, size int32,
) (*pb.FindRandomFavoriteWordMeaningsResponse, error) {
	return service.client.FindRandomFavoriteWordMeanings(
		context.Background(),
		&pb.FindRandomFavoriteWordMeaningsRequest{
			UserId: userId,
			Size:   size,
		},
	)
}
