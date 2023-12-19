package service

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/database"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExamService interface {
	CreateExam(topic, description string, isPublic bool, userId string) (string, error)
}

type examService struct {
	logger      log.Logger
	errorLogger log.Logger
}

func New(logger log.Logger) ExamService {
	var examService ExamService = examService{logger, level.Error(logger)}
	examService = loggingMiddleware{logger, examService}
	return examService
}

func (examService examService) CreateExam(
	topic, description string, isPublic bool, userId string) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	collection := database.GetCollection("exams")
	exam := model.Exam{
		Topic:       topic,
		Description: description,
		IsPublic:    isPublic,
		Tags:        []string{},
		UserId:      userId,
	}
	result, err := collection.InsertOne(context.TODO(), exam)

	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf("CreateExam fail! error: %w", err)
	}

	examId := result.InsertedID.(primitive.ObjectID).Hex()
	logger.Log("examId", examId)
	return examId, nil
}
