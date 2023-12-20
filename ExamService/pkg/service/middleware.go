package service

import (
	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type loggingMiddleware struct {
	logger log.Logger
	next   ExamService
}

func (mw loggingMiddleware) CreateExam(
	topic, description string, isPublic bool, userId string,
) (examId string, err error) {
	defer func() {
		mw.logger.Log(
			"method", "CreateExam",
			"topic", topic,
			"description", description,
			"isPublic", isPublic,
			"userId", userId,
			"err", err)
	}()
	return mw.next.CreateExam(topic, description, isPublic, userId)
}

func (mw loggingMiddleware) UpdateExam(
	examId,
	topic, description string, isPublic bool, userId string,
) (updatedExamId string, err error) {
	defer func() {
		mw.logger.Log(
			"method", "UpdateExam",
			"examId", examId,
			"topic", topic,
			"description", description,
			"isPublic", isPublic,
			"userId", userId,
			"err", err)
	}()
	return mw.next.UpdateExam(examId, topic, description, isPublic, userId)
}

func (mw loggingMiddleware) FindExams(
	pageIndex, pageSize int64, userId string,
) (total, pageCount int64, exams []model.Exam, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindExams",
			"pageIndex", pageIndex,
			"pageSize", pageSize,
			"userId", userId,
			"err", err)
	}()
	return mw.next.FindExams(pageIndex, pageSize, userId)
}
