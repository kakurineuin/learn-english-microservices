package microservice

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
)

const WORD_SERVER_ADDRESS = "localhost:8090"

var (
	connections       = []*grpc.ClientConn{}
	wordServiceClient pb.WordServiceClient
)

func Connect() error {
	conn, err := grpc.Dial(
		WORD_SERVER_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	connections = append(connections, conn)
	wordServiceClient = pb.NewWordServiceClient(conn)

	return nil
}

func Disconnect() error {
	for _, conn := range connections {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func FindWordByDictionary(word, userId string) (*pb.FindWordByDictionaryResponse, error) {
	return wordServiceClient.FindWordByDictionary(
		context.Background(),
		&pb.FindWordByDictionaryRequest{
			Word:   word,
			UserId: userId,
		},
	)
}

func CreateFavoriteWordMeaning(
	userId, wordMeaningId string,
) (*pb.CreateFavoriteWordMeaningResponse, error) {
	return wordServiceClient.CreateFavoriteWordMeaning(
		context.Background(),
		&pb.CreateFavoriteWordMeaningRequest{
			UserId:        userId,
			WordMeaningId: wordMeaningId,
		},
	)
}
