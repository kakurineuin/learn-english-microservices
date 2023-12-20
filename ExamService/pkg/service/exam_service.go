package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/database"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type ExamService interface {
	CreateExam(topic, description string, isPublic bool, userId string) (string, error)
	UpdateExam(examId, topic, description string, isPublic bool, userId string) (string, error)
	FindExams(
		pageIndex, pageSize int64,
		userId string,
	) (total, pageCount int64, exams []model.Exam, err error)
	DeleteExam(examId, userId string) error
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
	topic, description string, isPublic bool, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	now := time.Now()
	collection := database.GetCollection("exams")
	exam := model.Exam{
		Topic:       topic,
		Description: description,
		IsPublic:    isPublic,
		Tags:        []string{},
		UserId:      userId,
		CreatedAt:   now,
		UpdatedAt:   now,
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

func (examService examService) UpdateExam(
	examId, topic, description string, isPublic bool, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	collection := database.GetCollection("exams")
	id, _ := primitive.ObjectIDFromHex(examId)
	filter := bson.D{
		{"_id", id},
		{"userId", userId},
	}
	update := bson.D{{"$set", bson.D{
		{"topic", topic},
		{"description", description},
		{"isPublic", isPublic},
		{"updatedAt", time.Now()},
	}}}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		errorLogger.Log("err", err)
		return "", err
	}

	// 查無符合條件的測驗可供修改
	if result.MatchedCount == 0 {
		err = fmt.Errorf("Exam not found by examId: %s, userId: %s", examId, userId)
		errorLogger.Log("err", err)
		return "", err
	}

	logger.Log("examId", examId)
	return examId, nil
}

func (examService examService) FindExams(
	pageIndex, pageSize int64,
	userId string,
) (total, pageCount int64, exams []model.Exam, err error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	collection := database.GetCollection("exams")
	filter := bson.D{{"userId", userId}}
	sort := bson.D{{"updatedAt", -1}} // descending
	opts := options.Find().SetSort(sort).SetSkip(pageSize * pageIndex).SetLimit(pageSize)
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, err
	}

	if err = cursor.All(context.TODO(), &exams); err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, err
	}

	// Total
	total, err = collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, err
	}

	// PageCount
	pageCount = int64(math.Ceil(float64(total) / float64(pageSize)))
	logger.Log("total", total, "pageCount", pageCount, "exams size", len(exams))
	return
}

func (examService examService) DeleteExam(examId, userId string) error {
	errorLogger := examService.errorLogger

	// TODO: 改用交易的方式同時刪除 Exam, Question, AnswerWrong, ExamRecord
	// https://www.mongodb.com/docs/drivers/go/current/fundamentals/transactions/

	collection := database.GetCollection("exams")
	id, _ := primitive.ObjectIDFromHex(examId)
	filter := bson.D{{"_id", id}, {"userId", userId}}
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		errorLogger.Log("err", err)
		return err
	}

	// 查無符合條件的測驗可供刪除
	if result.DeletedCount == 0 {
		err = fmt.Errorf("Exam not found by examId: %s, userId: %s", examId, userId)
		errorLogger.Log("err", err)
		return err
	}

	return nil
}
