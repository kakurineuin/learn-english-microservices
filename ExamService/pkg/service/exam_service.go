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

var unauthorizedOperationError = fmt.Errorf("Unauthorized operation")

type ExamService interface {
	// Exam
	CreateExam(topic, description string, isPublic bool, userId string) (string, error)
	UpdateExam(examId, topic, description string, isPublic bool, userId string) (string, error)
	FindExams(
		pageIndex, pageSize int64,
		userId string,
	) (total, pageCount int64, exams []model.Exam, err error)
	DeleteExam(examId, userId string) error

	// Question
	CreateQuestion(examId, ask string, answers []string, userId string) (string, error)
	UpdateQuestion(questionId, ask string, answers []string, userId string) (string, error)
	FindQuestions(
		pageIndex, pageSize int64,
		examId, userId string,
	) (total, pageCount int64, questions []model.Question, err error)
	DeleteQuestion(questionId, userId string) error

	// ExamRecord
	CreateExamRecord(examId string, score int64, wrongQuestionIds []string, userId string) error
	FindExamRecords(
		pageIndex, pageSize int64,
		examId, userId string,
	) (total, pageCount int64, examRecords []model.ExamRecord, err error)
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
	errorMessage := "CreateExam failed: %w"

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
		return "", fmt.Errorf(errorMessage, err)
	}

	logger.Log("examId", examId)
	return examId, nil
}

func (examService examService) UpdateExam(
	examId, topic, description string, isPublic bool, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger
	errorMessage := "UpdateExam failed: %w"

	databaseRepository := examService.databaseRepository
	exam, err := databaseRepository.GetExamById(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	if exam == nil {
		err = fmt.Errorf("Exam not found by id: %s", examId)
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	// 檢查使用者是否是該測驗的擁有者
	if exam.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	exam.Topic = topic
	exam.Description = description
	exam.IsPublic = isPublic
	err = databaseRepository.UpdateExam(context.TODO(), *exam)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
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
	errorMessage := "FindExams failed: %w"

	skip := pageSize * pageIndex
	exams, err = examService.databaseRepository.FindExamsByUserIdOrderByUpdateAtDesc(
		context.TODO(), userId, skip, pageSize)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// Total
	total, err = examService.databaseRepository.CountExamsByUserId(context.TODO(), userId)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// PageCount
	pageCount = int64(math.Ceil(float64(total) / float64(pageSize)))
	logger.Log("total", total, "pageCount", pageCount, "exams size", len(exams))
	return
}

func (examService examService) DeleteExam(examId, userId string) error {
	errorLogger := examService.errorLogger
	errorMessage := "DeleteExam failed: %w"

	databaseRepository := examService.databaseRepository

	// 檢查使用者是否是該測驗的擁有者
	exam, err := databaseRepository.GetExamById(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	if exam == nil {
		err = fmt.Errorf("Exam not found by id: %s", examId)
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	// 使用者不是該測驗的擁有者
	if exam.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	_, err = examService.databaseRepository.WithTransaction(
		func(ctx context.Context) (interface{}, error) {
			// Delete Exam
			_, err := databaseRepository.DeleteExamById(ctx, examId)
			if err != nil {
				return nil, err
			}

			// Delete Question
			_, err = databaseRepository.DeleteQuestionsByExamId(ctx, examId)
			if err != nil {
				return nil, err
			}

			// Delete AnswerWrong
			_, err = databaseRepository.DeleteAnswerWrongsByExamId(ctx, examId)
			if err != nil {
				return nil, err
			}

			// Delete ExamRecord
			_, err = databaseRepository.DeleteExamRecordsByExamId(ctx, examId)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	return nil
}

func (examService examService) CreateQuestion(
	examId, ask string, answers []string, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger
	errorMessage := "CreateQuestion failed: %w"

	databaseRepository := examService.databaseRepository
	exam, err := databaseRepository.GetExamById(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	if exam == nil {
		err = fmt.Errorf("Exam not found by id: %s", examId)
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	// 檢查使用者是否是該測驗的擁有者
	if exam.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	questionId, err := databaseRepository.CreateQuestion(context.TODO(), model.Question{
		ExamId:  examId,
		Ask:     ask,
		Answers: answers,
		UserId:  userId,
	})
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	logger.Log("questionId", questionId)
	return questionId, nil
}

func (examService examService) UpdateQuestion(
	questionId, ask string, answers []string, userId string,
) (string, error) {
	logger := examService.logger
	errorLogger := examService.errorLogger
	errorMessage := "UpdateQuestion failed: %w"

	databaseRepository := examService.databaseRepository
	question, err := databaseRepository.GetQuestionById(context.TODO(), questionId)
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	if question == nil {
		err = fmt.Errorf("Question not found by id: %s", questionId)
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	// 檢查使用者是否是該 question 的擁有者
	if question.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
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
		_, err = databaseRepository.DeleteAnswerWrongsByQuestionId(ctx, questionId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		errorLogger.Log("err", err)
		return "", fmt.Errorf(errorMessage, err)
	}

	logger.Log("questionId", questionId)
	return questionId, nil
}

func (examService examService) FindQuestions(
	pageIndex, pageSize int64,
	examId, userId string,
) (total, pageCount int64, questions []model.Question, err error) {
	logger := examService.logger
	errorLogger := examService.errorLogger
	errorMessage := "FindQuestions failed: %w"

	databaseRepository := examService.databaseRepository

	exam, err := databaseRepository.GetExamById(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	if exam == nil {
		err := fmt.Errorf("Exam not found by id: %s", examId)
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// 檢查不能查詢別人的 question
	if exam.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	skip := pageSize * pageIndex
	questions, err = databaseRepository.FindQuestionsByExamIdOrderByUpdateAtDesc(
		context.TODO(), examId, skip, pageSize)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// Total
	total, err = databaseRepository.CountQuestionsByExamId(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// PageCount
	pageCount = int64(math.Ceil(float64(total) / float64(pageSize)))
	logger.Log("total", total, "pageCount", pageCount, "questions size", len(questions))
	return
}

func (examService examService) DeleteQuestion(
	questionId, userId string,
) error {
	errorLogger := examService.errorLogger
	errorMessage := "DeleteQuestion failed: %w"

	databaseRepository := examService.databaseRepository

	// 檢查不能刪除別人的 question
	question, err := databaseRepository.GetQuestionById(context.TODO(), questionId)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	if question == nil {
		err = fmt.Errorf("Question not found by id: %s", questionId)
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	if question.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	_, err = databaseRepository.WithTransaction(func(ctx context.Context) (interface{}, error) {
		// Delete AnswerWrong
		_, err = databaseRepository.DeleteAnswerWrongsByQuestionId(ctx, questionId)
		if err != nil {
			return nil, err
		}

		// Delete Question
		_, err = databaseRepository.DeleteQuestionById(ctx, questionId)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	return nil
}

func (examService examService) CreateExamRecord(
	examId string, score int64, wrongQuestionIds []string, userId string,
) error {
	errorLogger := examService.errorLogger
	errorMessage := "CreateExamRecord failed: %w"

	databaseRepository := examService.databaseRepository
	exam, err := databaseRepository.GetExamById(context.TODO(), examId)
	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	if exam == nil {
		err := fmt.Errorf("Exam not found by id: %s", examId)
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	// 檢查不能新增別人的測驗紀錄
	if !exam.IsPublic && exam.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	_, err = databaseRepository.WithTransaction(func(ctx context.Context) (interface{}, error) {
		// 新增測驗紀錄
		_, err = databaseRepository.CreateExamRecord(ctx, model.ExamRecord{
			ExamId: examId,
			Score:  score,
			UserId: userId,
		})
		if err != nil {
			return nil, err
		}

		// 更新問題的答錯次數
		for _, questionId := range wrongQuestionIds {
			_, _, err := databaseRepository.UpsertAnswerWrongByTimesPlusOne(
				ctx,
				examId,
				questionId,
				userId,
			)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		errorLogger.Log("err", err)
		return fmt.Errorf(errorMessage, err)
	}

	return nil
}

func (examService examService) FindExamRecords(
	pageIndex, pageSize int64,
	examId, userId string,
) (total, pageCount int64, examRecords []model.ExamRecord, err error) {
	logger := examService.logger
	errorLogger := examService.errorLogger
	errorMessage := "FindExamRecords failed: %w"

	databaseRepository := examService.databaseRepository
	skip := pageSize * pageIndex
	examRecords, err = databaseRepository.FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
		context.TODO(), examId, userId, skip, pageSize)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// Total
	total, err = databaseRepository.CountExamRecordsByExamIdAndUserId(
		context.TODO(),
		examId,
		userId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// PageCount
	pageCount = int64(math.Ceil(float64(total) / float64(pageSize)))
	logger.Log("total", total, "pageCount", pageCount, "examRecords size", len(examRecords))
	return
}
