package service

import (
	"context"
	"time"

	"github.com/go-kit/log"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type loggingMiddleware struct {
	logger log.Logger
	next   ExamService
}

func (mw loggingMiddleware) CreateExam(
	ctx context.Context, topic, description string, isPublic bool, userId string,
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
	return mw.next.CreateExam(ctx, topic, description, isPublic, userId)
}

func (mw loggingMiddleware) UpdateExam(
	ctx context.Context,
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
	return mw.next.UpdateExam(ctx, examId, topic, description, isPublic, userId)
}

func (mw loggingMiddleware) FindExams(
	ctx context.Context, pageIndex, pageSize int32, userId string,
) (total, pageCount int32, exams []model.Exam, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindExams",
			"pageIndex", pageIndex,
			"pageSize", pageSize,
			"userId", userId,
			"err", err)
	}()
	return mw.next.FindExams(ctx, pageIndex, pageSize, userId)
}

func (mw loggingMiddleware) DeleteExam(
	ctx context.Context, examId, userId string,
) (err error) {
	defer func() {
		mw.logger.Log(
			"method", "DeleteExam",
			"examId", examId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.DeleteExam(ctx, examId, userId)
}

func (mw loggingMiddleware) CreateQuestion(
	ctx context.Context, examId, ask string, answers []string, userId string,
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
	return mw.next.CreateQuestion(ctx, examId, ask, answers, userId)
}

func (mw loggingMiddleware) UpdateQuestion(
	ctx context.Context, questionId, ask string, answers []string, userId string,
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
	return mw.next.UpdateQuestion(ctx, questionId, ask, answers, userId)
}

func (mw loggingMiddleware) FindQuestions(
	ctx context.Context, pageIndex, pageSize int32, examId, userId string,
) (total, pageCount int32, questions []model.Question, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindQuestions",
			"pageIndex", pageIndex,
			"pageSize", pageSize,
			"exmaId", examId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.FindQuestions(ctx, pageIndex, pageSize, examId, userId)
}

func (mw loggingMiddleware) DeleteQuestion(
	ctx context.Context, questionId, userId string,
) (err error) {
	defer func() {
		mw.logger.Log(
			"method", "DeleteQuestion",
			"questionId", questionId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.DeleteQuestion(ctx, questionId, userId)
}

func (mw loggingMiddleware) CreateExamRecord(
	ctx context.Context, examId string, score int32, wrongQuestionIds []string, userId string,
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
	return mw.next.CreateExamRecord(ctx, examId, score, wrongQuestionIds, userId)
}

func (mw loggingMiddleware) FindExamRecords(
	ctx context.Context, pageIndex, pageSize int32, examId, userId string,
) (total, pageCount int32, examRecords []model.ExamRecord, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindExamRecords",
			"pageIndex", pageIndex,
			"pageSize", pageSize,
			"exmaId", examId,
			"userId", userId,
			"err", err)
	}()
	return mw.next.FindExamRecords(ctx, pageIndex, pageSize, examId, userId)
}

func (mw loggingMiddleware) FindExamRecordOverview(
	ctx context.Context, examId, userId string, startDate time.Time,
) (
	strStartDate string,
	exam *model.Exam,
	questions []model.Question,
	answerWrongs []model.AnswerWrong,
	examRecords []model.ExamRecord,
	err error,
) {
	defer func() {
		mw.logger.Log(
			"method", "FindExamRecordOverview",
			"exmaId", examId,
			"userId", userId,
			"startDate", startDate,
			"err", err)
	}()
	return mw.next.FindExamRecordOverview(ctx, examId, userId, startDate)
}

func (mw loggingMiddleware) FindExamInfos(
	ctx context.Context, userId string, isPublic bool,
) (examInfos []ExamInfo, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindExamInfos",
			"userId", userId,
			"isPublic", isPublic,
			"err", err)
	}()
	return mw.next.FindExamInfos(ctx, userId, isPublic)
}

func (mw loggingMiddleware) FindRandomQuestions(
	ctx context.Context, examId, userId string, size int32,
) (exam *model.Exam, questions []model.Question, err error) {
	defer func() {
		mw.logger.Log(
			"method", "FindRandomQuestions",
			"examId", examId,
			"userId", userId,
			"size", size,
			"err", err)
	}()
	return mw.next.FindRandomQuestions(ctx, examId, userId, size)
}
