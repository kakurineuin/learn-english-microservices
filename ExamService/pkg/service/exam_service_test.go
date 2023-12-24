package service

import (
	"log"
	"os"
	"testing"

	gokitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/repository"
)

const (
	EXAM_ID     = "6585ebdb460ab3291a0e555f"
	QUESTION_ID = "6585fc6430ca10723a6453b8"
)

type MyTestSuite struct {
	suite.Suite
	examService            ExamService
	mockDatabaseRepository *repository.MockDatabaseRepository
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}

// run once, before test suite methods
func (s *MyTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	logger := gokitlog.NewJSONLogger(os.Stdout)
	logger = gokitlog.With(
		logger,
		"ts",
		gokitlog.DefaultTimestampUTC,
		"caller",
		gokitlog.DefaultCaller,
	)
	mockDatabaseRepository := repository.NewMockDatabaseRepository(s.T())
	s.examService = New(logger, mockDatabaseRepository)
	s.mockDatabaseRepository = mockDatabaseRepository
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
}

// run before each test
func (s *MyTestSuite) SetupTest() {
	log.Println("SetupTest()")
}

// run after each test
func (s *MyTestSuite) TearDownTest() {
	log.Println("TearDownTest()")
}

// run before each test
func (s *MyTestSuite) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
}

// run after each test
func (s *MyTestSuite) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}

func (s *MyTestSuite) TestCreateExam() {
	s.mockDatabaseRepository.EXPECT().CreateExam(mock.Anything, mock.Anything).Return("exam01", nil)

	examId, err := s.examService.CreateExam(
		"topic01",
		"desc01",
		true,
		"userId",
	)
	s.Nil(err)
	s.NotEmpty(examId)
}

func (s *MyTestSuite) TestUpdateExam() {
	userId := "user01"
	id, err := primitive.ObjectIDFromHex(EXAM_ID)

	s.mockDatabaseRepository.EXPECT().GetExamById(mock.Anything, EXAM_ID).Return(&model.Exam{
		Id:          id,
		Topic:       "topic01",
		Description: "desc01",
		Tags:        []string{"t01"},
		IsPublic:    true,
		UserId:      userId,
	}, nil)
	s.mockDatabaseRepository.EXPECT().UpdateExam(mock.Anything, mock.Anything).Return(nil)

	examId, err := s.examService.UpdateExam(
		EXAM_ID,
		"topic02",
		"desc02",
		false,
		userId,
	)
	s.Nil(err)
	s.Equal(EXAM_ID, examId)
}

func (s *MyTestSuite) TestFindExams() {
	userId := "user01"

	var pageIndex int64 = 0
	var pageSize int64 = 10
	s.mockDatabaseRepository.EXPECT().FindExamsByUserIdOrderByUpdateAtDesc(
		mock.Anything,
		userId,
		pageIndex*pageSize,
		pageSize,
	).Return([]model.Exam{
		{},
		{},
		{},
	}, nil)

	var expectedTotal int64 = 3
	s.mockDatabaseRepository.EXPECT().
		CountExamsByUserId(mock.Anything, userId).
		Return(expectedTotal, nil)

	total, pageCount, exams, err := s.examService.FindExams(pageIndex, pageSize, userId)
	s.Nil(err)
	s.Equal(expectedTotal, total)
	s.Equal(int64(1), pageCount)
	s.Equal(3, len(exams))
}

func (s *MyTestSuite) TestDeleteExam() {
	userId := "user01"

	id, err := primitive.ObjectIDFromHex(EXAM_ID)
	s.Nil(err)

	s.mockDatabaseRepository.EXPECT().GetExamById(mock.Anything, EXAM_ID).Return(&model.Exam{
		Id:          id,
		Topic:       "topic01",
		Description: "desc01",
		Tags:        []string{"t01"},
		IsPublic:    true,
		UserId:      userId,
	}, nil)
	s.mockDatabaseRepository.EXPECT().WithTransaction(mock.Anything).Return(nil, nil)

	err = s.examService.DeleteExam(EXAM_ID, userId)
	s.Nil(err)
}

func (s *MyTestSuite) TestCreateQuestion() {
	userId := "user01"
	expectedQuestionId := "q01"

	s.mockDatabaseRepository.EXPECT().
		CreateQuestion(mock.Anything, mock.Anything).
		Return(expectedQuestionId, nil)

	questionId, err := s.examService.CreateQuestion(
		EXAM_ID,
		"ask",
		[]string{"a01", "a02"},
		userId)
	s.Nil(err)
	s.Equal(expectedQuestionId, questionId)
}

func (s *MyTestSuite) TestUpdaetQuestion() {
	userId := "user01"

	id, err := primitive.ObjectIDFromHex(QUESTION_ID)
	s.Nil(err)

	s.mockDatabaseRepository.EXPECT().
		GetQuestionById(mock.Anything, QUESTION_ID).
		Return(&model.Question{
			Id:      id,
			ExamId:  EXAM_ID,
			Ask:     "ask01",
			Answers: []string{"a01", "a02"},
			UserId:  userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		WithTransaction(mock.Anything).
		Return(nil, nil)

	questionId, err := s.examService.UpdateQuestion(
		QUESTION_ID,
		"ask02",
		[]string{"b01", "b02"},
		userId)
	s.Nil(err)
	s.Equal(QUESTION_ID, questionId)
}

func (s *MyTestSuite) TestFindQuestions() {
	userId := "user01"
	var pageIndex int64 = 2
	var pageSize int64 = 10
	skip := pageIndex * pageSize
	limit := pageSize
	var expectedTotal int64 = 23

	s.mockDatabaseRepository.EXPECT().
		FindQuestionsByExamIdOrderByUpdateAtDesc(mock.Anything, EXAM_ID, skip, limit).
		Return([]model.Question{
			{},
			{},
			{},
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		CountQuestionsByExamIdAndUserId(mock.Anything, EXAM_ID, userId).
		Return(expectedTotal, nil)

	total, pageCount, questions, err := s.examService.FindQuestions(
		pageIndex,
		pageSize,
		EXAM_ID,
		userId,
	)
	s.Nil(err)
	s.Equal(expectedTotal, total)
	s.Equal(int64(3), pageCount)
	s.Equal(3, len(questions))
}

func (s *MyTestSuite) TestDeleteQuestion() {
	userId := "user01"
	id, err := primitive.ObjectIDFromHex(QUESTION_ID)
	s.Nil(err)

	s.mockDatabaseRepository.EXPECT().
		GetQuestionById(mock.Anything, QUESTION_ID).
		Return(&model.Question{
			Id:      id,
			ExamId:  EXAM_ID,
			Ask:     "ask02",
			Answers: []string{"b01", "b02"},
			UserId:  userId,
		}, nil)

	err = s.examService.DeleteQuestion(QUESTION_ID, userId)
	s.Nil(err)
}

func (s *MyTestSuite) TestCreateExamRecord() {
	userId := "user01"
	id, err := primitive.ObjectIDFromHex(EXAM_ID)
	s.Nil(err)

	s.mockDatabaseRepository.EXPECT().GetExamById(mock.Anything, EXAM_ID).Return(&model.Exam{
		Id:          id,
		Topic:       "topic",
		Description: "desc",
		Tags:        []string{"a01"},
		IsPublic:    true,
		UserId:      userId,
	}, nil)
	s.mockDatabaseRepository.EXPECT().WithTransaction(mock.Anything).Return(nil, nil)

	err = s.examService.CreateExamRecord(
		EXAM_ID,
		10,
		[]string{"question01", "question02", "quesiton03"},
		userId,
	)
	s.Nil(err)
}
