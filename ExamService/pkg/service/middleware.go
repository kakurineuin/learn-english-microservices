package service

import (
	"github.com/go-kit/log"
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
