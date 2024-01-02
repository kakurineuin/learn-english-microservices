package examservice

import (
	"context"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
)

//go:generate mockery --name ExamService
type ExamService interface {
	Connect() error
	Disconnect() error
	CreateExam(
		topic, description string, isPublic bool, userId string,
	) (*pb.CreateExamResponse, error)
	FindExams(
		pageIndex, pageSize int32, userId string,
	) (*pb.FindExamsResponse, error)
	UpdateExam(
		examId, topic, description string, isPublic bool, userId string,
	) (*pb.UpdateExamResponse, error)
	DeleteExam(
		examId, userId string,
	) (*pb.DeleteExamResponse, error)
	FindQuestions(
		pageIndex, pageSize int32, examId, userId string,
	) (*pb.FindQuestionsResponse, error)
	CreateQuestion(
		examId, ask string, answers []string, userId string,
	) (*pb.CreateQuestionResponse, error)
	UpdateQuestion(
		questionId, ask string, answers []string, userId string,
	) (*pb.UpdateQuestionResponse, error)
	DeleteQuestion(
		questionId, userId string,
	) (*pb.DeleteQuestionResponse, error)
}

func New(serverAddress string) ExamService {
	return &examService{
		serverAddress: serverAddress,
	}
}

type examService struct {
	serverAddress string
	connection    *grpc.ClientConn
	client        pb.ExamServiceClient
}

func (service *examService) Connect() error {
	log.Infof("Start to connect ExamService at %s", service.serverAddress)

	conn, err := grpc.Dial(
		service.serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}

	service.connection = conn
	service.client = pb.NewExamServiceClient(conn)

	return nil
}

func (service *examService) Disconnect() error {
	if err := service.connection.Close(); err != nil {
		return err
	}

	service.connection = nil
	service.client = nil

	return nil
}

func (service examService) CreateExam(
	topic,
	description string,
	isPublic bool,
	userId string,
) (*pb.CreateExamResponse, error) {
	return service.client.CreateExam(
		context.Background(),
		&pb.CreateExamRequest{
			Topic:       topic,
			Description: description,
			IsPublic:    isPublic,
			UserId:      userId,
		},
	)
}

func (service examService) FindExams(
	pageIndex, pageSize int32,
	userId string,
) (*pb.FindExamsResponse, error) {
	return service.client.FindExams(
		context.Background(),
		&pb.FindExamsRequest{
			PageIndex: pageIndex,
			PageSize:  pageSize,
			UserId:    userId,
		},
	)
}

func (examServiceClient examService) UpdateExam(
	examId,
	topic,
	description string,
	isPublic bool,
	userId string,
) (*pb.UpdateExamResponse, error) {
	return examServiceClient.client.UpdateExam(
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

func (examServiceClient examService) DeleteExam(
	examId, userId string,
) (*pb.DeleteExamResponse, error) {
	return examServiceClient.client.DeleteExam(
		context.Background(),
		&pb.DeleteExamRequest{
			ExamId: examId,
			UserId: userId,
		},
	)
}

func (service examService) FindQuestions(
	pageIndex, pageSize int32,
	examId, userId string,
) (*pb.FindQuestionsResponse, error) {
	return service.client.FindQuestions(
		context.Background(),
		&pb.FindQuestionsRequest{
			PageIndex: pageIndex,
			PageSize:  pageSize,
			ExamId:    examId,
			UserId:    userId,
		},
	)
}

func (service examService) CreateQuestion(
	examId, ask string, answers []string, userId string,
) (*pb.CreateQuestionResponse, error) {
	return service.client.CreateQuestion(
		context.Background(),
		&pb.CreateQuestionRequest{
			ExamId:  examId,
			Ask:     ask,
			Answers: answers,
			UserId:  userId,
		},
	)
}

func (service examService) UpdateQuestion(
	questionId, ask string, answers []string, userId string,
) (*pb.UpdateQuestionResponse, error) {
	return service.client.UpdateQuestion(
		context.Background(),
		&pb.UpdateQuestionRequest{
			QuestionId: questionId,
			Ask:        ask,
			Answers:    answers,
			UserId:     userId,
		},
	)
}

func (service examService) DeleteQuestion(
	questionId, userId string,
) (*pb.DeleteQuestionResponse, error) {
	return service.client.DeleteQuestion(
		context.Background(),
		&pb.DeleteQuestionRequest{
			QuestionId: questionId,
			UserId:     userId,
		},
	)
}
