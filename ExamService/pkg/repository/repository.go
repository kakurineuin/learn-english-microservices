package repository

import (
	"context"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type transactionFunc func(ctx context.Context) (interface{}, error)

type DatabaseRepository interface {
	ConnectDB(uri string) error
	DisconnectDB() error

	// Transaction
	WithTransaction(transactoinFunc transactionFunc) (interface{}, error)

	// Exam
	CreateExam(ctx context.Context, exam model.Exam) (examId string, err error)
	UpdateExam(ctx context.Context, exam model.Exam) error
	GetExam(ctx context.Context, examId string) (exam *model.Exam, err error)
	FindExamsOrderByUpdateAtDesc(
		ctx context.Context,
		userId string,
		skip, limit int64,
	) (exams []model.Exam, err error)
	DeleteExam(ctx context.Context, examId string) error

	// Question
	// CreateQuestion(ctx context.Context, question model.Question) (questionId string, err error)
	// UpdateQuestion(ctx context.Context, question model.Question) error

	// TODO: AnswerWrong
}
