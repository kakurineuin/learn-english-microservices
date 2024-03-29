package repository

import (
	"context"
	"time"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type transactionFunc func(ctx context.Context) (interface{}, error)

//go:generate mockery --name DatabaseRepository
type DatabaseRepository interface {
	ConnectDB(ctx context.Context, uri string) error
	DisconnectDB(ctx context.Context) error

	// Transaction
	WithTransaction(ctx context.Context, transactoinFunc transactionFunc) (interface{}, error)

	// Exam
	CreateExam(ctx context.Context, exam model.Exam) (examId string, err error)
	UpdateExam(ctx context.Context, exam model.Exam) error
	GetExamById(ctx context.Context, examId string) (exam *model.Exam, err error)
	FindExamsByUserIdOrderByUpdateAtDesc(
		ctx context.Context,
		userId string,
		skip, limit int32,
	) (exams []model.Exam, err error)
	FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(
		ctx context.Context,
		userId string,
		isPublic bool,
	) (exams []model.Exam, err error)
	DeleteExamById(ctx context.Context, examId string) (deletedCount int32, err error)
	CountExamsByUserId(ctx context.Context, userId string) (count int32, err error)

	// Question
	CreateQuestion(ctx context.Context, question model.Question) (questionId string, err error)
	UpdateQuestion(ctx context.Context, question model.Question) error
	GetQuestionById(ctx context.Context, questionId string) (question *model.Question, err error)
	FindQuestionsByQuestionIds(
		ctx context.Context,
		questionIds []string,
	) (questions []model.Question, err error)
	FindQuestionsByExamIdOrderByUpdateAtDesc(
		ctx context.Context,
		examId string,
		skip, limit int32,
	) (questions []model.Question, err error)
	DeleteQuestionById(ctx context.Context, questionId string) (deletedCount int32, err error)
	DeleteQuestionsByExamId(ctx context.Context, examId string) (deletedCount int32, err error)
	CountQuestionsByExamId(
		ctx context.Context,
		examId string,
	) (count int32, err error)

	// AnswerWrong
	DeleteAnswerWrongsByQuestionId(
		ctx context.Context,
		questionId string,
	) (deletedCount int32, err error)
	DeleteAnswerWrongsByExamId(ctx context.Context, examId string) (deletedCount int32, err error)
	UpsertAnswerWrongByTimesPlusOne(
		ctx context.Context,
		examId, questionId, userId string,
	) (modifiedCount, upsertedCount int32, err error)
	FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc(
		ctx context.Context,
		examId, userId string,
		limit int32,
	) (answerWrongs []model.AnswerWrong, err error)

	// ExamRecord
	CreateExamRecord(
		ctx context.Context,
		examRecord model.ExamRecord,
	) (examRecordId string, err error)
	DeleteExamRecordsByExamId(ctx context.Context, examId string) (deletedCount int32, err error)
	FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
		ctx context.Context,
		examId, userId string,
		skip, limit int32,
	) (examRecords []model.ExamRecord, err error)
	CountExamRecordsByExamIdAndUserId(
		ctx context.Context,
		examId, userId string,
	) (count int32, err error)
	FindExamRecordsByExamIdAndUserIdAndCreatedAt(
		ctx context.Context,
		examId,
		userId string,
		createdAt time.Time,
	) (examRecords []model.ExamRecord, err error)
}
