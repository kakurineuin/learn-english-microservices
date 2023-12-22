package service

import (
	"context"
	"fmt"
	"math"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/repository"
)

// TODO 全部方法改用 repository 去對資料庫存取

type ExamService interface {
	CreateExam(topic, description string, isPublic bool, userId string) (string, error)
	UpdateExam(examId, topic, description string, isPublic bool, userId string) (string, error)
	FindExams(
		pageIndex, pageSize int64,
		userId string,
	) (total, pageCount int64, exams []model.Exam, err error)
	DeleteExam(examId, userId string) error

	CreateQuestion(examId, ask string, answers []string, userId string) (string, error)
	UpdateQuestion(questionId, ask string, answers []string, userId string) (string, error)
}

type examService struct {
	logger             log.Logger
	errorLogger        log.Logger
	databaseRepository repository.DatabaseRepository
}

func New(logger log.Logger, databaseRepository repository.DatabaseRepository) ExamService {
	var examService ExamService = examService{
		logger:             logger,
		errorLogger:        level.Error(logger),
		databaseRepository: databaseRepository,
	}
	examService = loggingMiddleware{logger, examService}
	return examService
}

func (examService examService) CreateExam(
	topic, description string, isPublic bool, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	exam := model.Exam{
		Topic:       topic,
		Description: description,
		IsPublic:    isPublic,
		Tags:        []string{},
		UserId:      userId,
	}
	examId, err := examService.databaseRepository.CreateExam(context.TODO(), exam)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf("CreateExam fail! error: %w", err)
	}

	logger.Log("examId", examId)
	return examId, nil
}

func (examService examService) UpdateExam(
	examId, topic, description string, isPublic bool, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	databaseRepository := examService.databaseRepository
	exam, err := databaseRepository.GetExam(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return "", err
	}

	// 檢查使用者是否是該測驗的擁有者
	if exam.UserId != userId {
		err = fmt.Errorf("Unauthorized operation")
		errorLogger.Log("err", err)
		return "", err
	}

	exam.Topic = topic
	exam.Description = description
	exam.IsPublic = isPublic
	err = examService.databaseRepository.UpdateExam(context.TODO(), *exam)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf("UpdateExam fail! error: %w", err)
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

	skip := pageSize * pageIndex
	exams, err = examService.databaseRepository.FindExamsOrderByUpdateAtDesc(
		context.TODO(), userId, skip, pageSize)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf("FindExams fail! error: %w", err)
	}

	// Total
	total, err = examService.databaseRepository.CountExamsByUserId(context.TODO(), userId)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf("FindExams fail! error: %w", err)
	}

	// PageCount
	pageCount = int64(math.Ceil(float64(total) / float64(pageSize)))
	logger.Log("total", total, "pageCount", pageCount, "exams size", len(exams))
	return
}

func (examService examService) DeleteExam(examId, userId string) error {
	errorLogger := examService.errorLogger
	databaseRepository := examService.databaseRepository

	// 檢查使用者是否是該測驗的擁有者
	isExist, err := databaseRepository.ExistExam(context.TODO(), examId, userId)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf("DeleteExam fail! error: %w", err)
	}

	// 使用者不是該測驗的擁有者
	if !isExist {
		err = fmt.Errorf("Unauthorized operation")
		errorLogger.Log("err", err)
		return err
	}

	_, err = examService.databaseRepository.WithTransaction(
		func(ctx context.Context) (interface{}, error) {
			// Delete Exam
			err := databaseRepository.DeleteExam(ctx, examId)
			if err != nil {
				return nil, err
			}

			// Delete Question
			err = databaseRepository.DeleteQuestionsByExamId(ctx, examId)

			// TODO: Delete AnswerWrong

			// TODO: Delete ExamRecord

			return nil, nil
		},
	)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf("DeleteExam fail! error: %w", err)
	}

	return nil
}

func (examService examService) CreateQuestion(
	examId, ask string, answers []string, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	questionId, err := examService.databaseRepository.CreateQuestion(context.TODO(), model.Question{
		ExamId:  examId,
		Ask:     ask,
		Answers: answers,
		UserId:  userId,
	})
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf("CreateQuestion fail! error: %w", err)
	}

	logger.Log("questionId", questionId)
	return questionId, nil
}

func (examService examService) UpdateQuestion(
	questionId, ask string, answers []string, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger

	databaseRepository := examService.databaseRepository
	question, err := databaseRepository.GetQuestion(context.TODO(), questionId)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf("UpdateQuestion fail! error: %w", err)
	}

	// 檢查使用者是否是該 question 的擁有者
	if question.UserId != userId {
		err = fmt.Errorf("Unauthorized operation")
		errorLogger.Log("err", err)
		return "", err
	}

	_, err = databaseRepository.WithTransaction(func(ctx context.Context) (interface{}, error) {
		// 修改 Question
		question.Ask = ask
		question.Answers = answers
		err = databaseRepository.UpdateQuestion(ctx, *question)
		if err != nil {
			return nil, err
		}

		// 刪除相關的 AnswerWrong
		err = databaseRepository.DeleteAnswerWrongByQuestionId(ctx, questionId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf("UpdateQuestion fail! error: %w", err)
	}

	logger.Log("questionId", questionId)
	return questionId, nil
}
