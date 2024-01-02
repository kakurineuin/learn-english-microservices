package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/repository"
)

var unauthorizedOperationError = fmt.Errorf("Unauthorized operation")

// 這不是 MongoDB 的 document，所以不放在 model 目錄之下
type ExamInfo struct {
	ExamId        string
	Topic         string
	Description   string
	IsPublic      bool
	QuestionCount int32
	RecordCount   int32
}

type ExamService interface {
	// Exam
	CreateExam(topic, description string, isPublic bool, userId string) (string, error)
	UpdateExam(examId, topic, description string, isPublic bool, userId string) (string, error)
	FindExams(
		pageIndex, pageSize int32,
		userId string,
	) (total, pageCount int32, exams []model.Exam, err error)
	DeleteExam(examId, userId string) error

	// Question
	CreateQuestion(examId, ask string, answers []string, userId string) (string, error)
	UpdateQuestion(questionId, ask string, answers []string, userId string) (string, error)
	FindQuestions(
		pageIndex, pageSize int32,
		examId, userId string,
	) (total, pageCount int32, questions []model.Question, err error)
	DeleteQuestion(questionId, userId string) error
	FindRandomQuestions(
		examId, userId string, size int32,
	) (exam *model.Exam, questions []model.Question, err error)

	// ExamRecord
	CreateExamRecord(examId string, score int32, wrongQuestionIds []string, userId string) error
	FindExamRecords(
		pageIndex, pageSize int32,
		examId, userId string,
	) (total, pageCount int32, examRecords []model.ExamRecord, err error)

	// ExamInfo
	FindExamInfos(userId string, isPublic bool) (examInfos []ExamInfo, err error)
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
	pageIndex, pageSize int32,
	userId string,
) (total, pageCount int32, exams []model.Exam, err error) {
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
	pageCount = int32(math.Ceil(float64(total) / float64(pageSize)))
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
	pageIndex, pageSize int32,
	examId, userId string,
) (total, pageCount int32, questions []model.Question, err error) {
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
	pageCount = int32(math.Ceil(float64(total) / float64(pageSize)))
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
	examId string, score int32, wrongQuestionIds []string, userId string,
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
	pageIndex, pageSize int32,
	examId, userId string,
) (total, pageCount int32, examRecords []model.ExamRecord, err error) {
	logger := examService.logger
	errorLogger := examService.errorLogger
	errorMessage := "FindExamRecords failed: %w"

	databaseRepository := examService.databaseRepository
	skip := pageSize * pageIndex
	limit := pageSize
	examRecords, err = databaseRepository.FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
		context.TODO(), examId, userId, skip, limit)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}

	// Total
	count, err := databaseRepository.CountExamRecordsByExamIdAndUserId(
		context.TODO(),
		examId,
		userId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return 0, 0, nil, fmt.Errorf(errorMessage, err)
	}
	total = int32(count)

	// PageCount
	pageCount = int32(math.Ceil(float64(total) / float64(pageSize)))
	logger.Log("total", total, "pageCount", pageCount, "examRecords size", len(examRecords))
	return
}

func (examService examService) FindExamInfos(
	userId string,
	isPublic bool,
) (examInfos []ExamInfo, err error) {
	errorLogger := examService.errorLogger
	errorMessage := "FindExamInfos failed: %w"

	databaseRepository := examService.databaseRepository

	exams, err := databaseRepository.FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(
		context.TODO(),
		userId,
		isPublic,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return nil, fmt.Errorf(errorMessage, err)
	}

	for _, exam := range exams {
		examId := exam.Id.Hex()
		questionCount, err := databaseRepository.CountQuestionsByExamId(context.TODO(), examId)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, fmt.Errorf(errorMessage, err)
		}

		if questionCount == 0 {
			continue
		}

		examRecordCount, err := databaseRepository.CountExamRecordsByExamIdAndUserId(
			context.TODO(),
			examId,
			userId,
		)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, fmt.Errorf(errorMessage, err)
		}

		examInfos = append(examInfos, ExamInfo{
			ExamId:        examId,
			Topic:         exam.Topic,
			Description:   exam.Description,
			IsPublic:      exam.IsPublic,
			QuestionCount: questionCount,
			RecordCount:   examRecordCount,
		})
	}

	return examInfos, nil
}

func (examService examService) FindRandomQuestions(
	examId, userId string, size int32,
) (exam *model.Exam, questions []model.Question, err error) {
	errorLogger := examService.errorLogger
	errorMessage := "FindRandomQuestions failed! error: %w"

	databaseRepository := examService.databaseRepository
	exam, err = databaseRepository.GetExamById(context.TODO(), examId)

	if err != nil {
		errorLogger.Log("err", err)
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	if exam == nil {
		return nil, []model.Question{}, nil
	}

	// 若測驗是不公開，則只有本人可以作測驗
	if !exam.IsPublic && exam.UserId != userId {
		err = unauthorizedOperationError
		errorLogger.Log("err", err)
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	total, err := databaseRepository.CountQuestionsByExamId(
		context.TODO(),
		examId,
	)
	if err != nil {
		errorLogger.Log("err", err)
		return nil, nil, fmt.Errorf(errorMessage, err)
	}

	// 總數是零，表示使用者還沒新增 Question
	if total == 0 {
		return exam, []model.Question{}, nil
	}

	indexes := []int32{}

	for i := int32(0); i < total; i++ {
		indexes = append(indexes, i)
	}

	// 將 indexes 隨機洗牌，然後根據洗牌後的 indexes 順序去查詢，達到隨機排序的效果
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})

	maxQueryTotal := min(total, size)

	for i := int32(0); i < maxQueryTotal; i++ {
		findQuestions, err := databaseRepository.FindQuestionsByExamIdOrderByUpdateAtDesc(
			context.TODO(),
			examId,
			indexes[i],
			1,
		)
		if err != nil {
			errorLogger.Log("err", err)
			return nil, nil, fmt.Errorf(errorMessage, err)
		}

		questions = append(questions, findQuestions[0])
	}

	return exam, questions, nil
}
