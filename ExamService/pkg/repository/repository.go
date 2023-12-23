package repository

import (
	"context"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type transactionFunc func(ctx context.Context) (interface{}, error)

//go:generate mockery --name DatabaseRepository
type DatabaseRepository interface {
	ConnectDB(uri string) error
	DisconnectDB() error

	// Transaction
	WithTransaction(transactoinFunc transactionFunc) (interface{}, error)

	// Exam
	CreateExam(ctx context.Context, exam model.Exam) (examId string, err error)
	UpdateExam(ctx context.Context, exam model.Exam) error
	GetExamById(ctx context.Context, examId string) (exam *model.Exam, err error)
	FindExamsByUserIdOrderByUpdateAtDesc(
		ctx context.Context,
		userId string,
		skip, limit int64,
	) (exams []model.Exam, err error)
	DeleteExamById(ctx context.Context, examId string) (deletedCount int64, err error)
	CountExamsByUserId(ctx context.Context, userId string) (count int64, err error)

	// Question
	CreateQuestion(ctx context.Context, question model.Question) (questionId string, err error)
	UpdateQuestion(ctx context.Context, question model.Question) error
	GetQuestionById(ctx context.Context, questionId string) (question *model.Question, err error)
	FindQuestionsByExamIdAndUserIdOrderByUpdateAtDesc(
		ctx context.Context,
		examId, userId string,
		skip, limit int64,
	) (questions []model.Question, err error)
	DeleteQuestionById(ctx context.Context, questionId string) (deletedCount int64, err error)
	DeleteQuestionsByExamId(ctx context.Context, examId string) (deletedCount int64, err error)
	CountQuestionsByExamIdAndUserId(
		ctx context.Context,
		examId, userId string,
	) (count int64, err error)

	// AnswerWrong
	DeleteAnswerWrongsByQuestionId(
		ctx context.Context,
		questionId string,
	) (deletedCount int64, err error)
	DeleteAnswerWrongsByExamId(ctx context.Context, examId string) (deletedCount int64, err error)
	UpsertAnswerWrongByTimesPlusOne(
		ctx context.Context,
		examId, questionId, userId string,
	) (modifiedCount, upsertedCount int64, err error)

	// ExamRecord
	CreateExamRecord(
		ctx context.Context,
		examRecord model.ExamRecord,
	) (examRecordId string, err error)
	DeleteExamRecordsByExamId(ctx context.Context, examId string) (deletedCount int64, err error)
}
