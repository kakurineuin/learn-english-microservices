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

func (mw loggingMiddleware) DeleteExam(
	examId, userId string,
) (err error) {
	defer func() {
		mw.logger.Log(
			"method", "DeleteExam",
			"examId", examId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.DeleteExam(examId, userId)
}

func (mw loggingMiddleware) CreateQuestion(
	examId, ask string, answers []string, userId string,
) (questionId string, err error) {
	defer func() {
		mw.logger.Log(
			"method", "CreateQuestion",
			"examId", examId,
			"ask", ask,
			"answers", answers,
			"userId", userId,
			"err", err)
	}()
	return mw.next.CreateQuestion(examId, ask, answers, userId)
}

func (mw loggingMiddleware) UpdateQuestion(
	questionId, ask string, answers []string, userId string,
) (updatedQuestionId string, err error) {
	defer func() {
		mw.logger.Log(
			"method", "UpdateQuestion",
			"questionId", questionId,
			"ask", ask,
			"answers", answers,
			"userId", userId,
			"err", err)
	}()
	return mw.next.UpdateQuestion(questionId, ask, answers, userId)
}

func (mw loggingMiddleware) FindQuestions(
	pageIndex, pageSize int64, examId, userId string,
) (total, pageCount int64, questions []model.Question, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindQuestions",
			"pageIndex", pageIndex,
			"pageSize", pageSize,
			"exmaId", examId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.FindQuestions(pageIndex, pageSize, examId, userId)
}

func (mw loggingMiddleware) DeleteQuestion(questionId, userId string) (err error) {
	defer func() {
		mw.logger.Log(
			"method", "DeleteQuestion",
			"questionId", questionId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.DeleteQuestion(questionId, userId)
}

func (mw loggingMiddleware) CreateExamRecord(
	examId string, score int64, wrongQuestionIds []string, userId string,
) (err error) {
	defer func() {
		mw.logger.Log(
			"method", "CreateExamRecord",
			"examId", examId,
			"score", score,
			"wrongQuestionIds", wrongQuestionIds,
			"userId", userId,
			"err", err)
	}()
	return mw.next.CreateExamRecord(examId, score, wrongQuestionIds, userId)
}

func (mw loggingMiddleware) FindExamRecords(
	pageIndex, pageSize int64, examId, userId string,
) (total, pageCount int64, examRecords []model.ExamRecord, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindExamRecords",
			"pageIndex", pageIndex,
			"pageSize", pageSize,
			"exmaId", examId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.FindExamRecords(pageIndex, pageSize, examId, userId)
}
