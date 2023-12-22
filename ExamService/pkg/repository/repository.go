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
	ExistExam(ctx context.Context, examId, userId string) (isExist bool, err error)
	CreateExam(ctx context.Context, exam model.Exam) (examId string, err error)
	UpdateExam(ctx context.Context, exam model.Exam) error
	GetExam(ctx context.Context, examId string) (exam *model.Exam, err error)
	FindExamsOrderByUpdateAtDesc(
		ctx context.Context,
		userId string,
		skip, limit int64,
	) (exams []model.Exam, err error)
	DeleteExam(ctx context.Context, examId string) error
	CountExamsByUserId(ctx context.Context, userId string) (count int64, err error)

	// Question
	CreateQuestion(ctx context.Context, question model.Question) (questionId string, err error)
	UpdateQuestion(ctx context.Context, question model.Question) error
	GetQuestion(ctx context.Context, questionId string) (question *model.Question, err error)
	FindQuestionsOrderByUpdateAtDesc(
		ctx context.Context,
		examId string,
		skip, limit int64,
	) (questions []model.Question, err error)
	DeleteQuestion(ctx context.Context, questionId string) error
	DeleteQuestionsByExamId(ctx context.Context, examId string) error

	// TODO: AnswerWrong
	DeleteAnswerWrongByQuestionId(ctx context.Context, questionId string) error
}
