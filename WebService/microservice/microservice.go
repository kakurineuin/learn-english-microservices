package microservice

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
)

const (
	EXAM_SERVER_ADDRESS = "localhost:8090"
	WORD_SERVER_ADDRESS = "localhost:8091"
)

var (
	connections       = []*grpc.ClientConn{}
	examServiceClient pb.ExamServiceClient
	wordServiceClient pb.WordServiceClient
)

func Connect() error {
	// ExamService
	conn, err := grpc.Dial(
		EXAM_SERVER_ADDRESS,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	connections = append(connections, conn)
	examServiceClient = pb.NewExamServiceClient(conn)

	// WordService
	conn, err = grpc.Dial(
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

func CreateExam(
	topic,
	description string,
	isPublic bool,
	userId string,
) (*pb.CreateExamResponse, error) {
	return examServiceClient.CreateExam(
		context.Background(),
		&pb.CreateExamRequest{
			Topic:       topic,
			Description: description,
			IsPublic:    isPublic,
			UserId:      userId,
		},
	)
}

func FindExams(
	pageIndex, pageSize int64,
	userId string,
) (*pb.FindExamsResponse, error) {
	return examServiceClient.FindExams(
		context.Background(),
		&pb.FindExamsRequest{
			PageIndex: pageIndex,
			PageSize:  pageSize,
			UserId:    userId,
		},
	)
}

func UpdateExam(
	examId,
	topic,
	description string,
	isPublic bool,
	userId string,
) (*pb.UpdateExamResponse, error) {
	return examServiceClient.UpdateExam(
		context.Background(),
		&pb.UpdateExamRequest{
			ExamId:      examId,
			Topic:       topic,
			Description: description,
			IsPublic:    isPublic,
			UserId:      userId,
		},
	)
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

func DeleteFavoriteWordMeaning(
	favoriteWordMeaningId, userId string,
) (*pb.DeleteFavoriteWordMeaningResponse, error) {
	return wordServiceClient.DeleteFavoriteWordMeaning(
		context.Background(),
		&pb.DeleteFavoriteWordMeaningRequest{
			FavoriteWordMeaningId: favoriteWordMeaningId,
			UserId:                userId,
		},
	)
}

func FindFavoriteWordMeanings(
	pageIndex, pageSize int64,
	userId, word string,
) (*pb.FindFavoriteWordMeaningsResponse, error) {
	return wordServiceClient.FindFavoriteWordMeanings(
		context.Background(),
		&pb.FindFavoriteWordMeaningsRequest{
			PageIndex: pageIndex,
			PageSize:  pageSize,
			UserId:    userId,
			Word:      word,
		},
	)
}

func FindRandomFavoriteWordMeanings(
	userId string, size int64,
) (*pb.FindRandomFavoriteWordMeaningsResponse, error) {
	return wordServiceClient.FindRandomFavoriteWordMeanings(
		context.Background(),
		&pb.FindRandomFavoriteWordMeaningsRequest{
			UserId: userId,
			Size:   size,
		},
	)
}
