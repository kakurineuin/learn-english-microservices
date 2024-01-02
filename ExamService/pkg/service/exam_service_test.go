package service

import (
	"log"
	"os"
	"testing"

	gokitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/repository"
)

type MyTestSuite struct {
	suite.Suite
	examService            examService
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
	s.examService = examService{
		logger:             logger,
		errorLogger:        level.Error(logger),
		databaseRepository: nil,
	}
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
}

// run before each test
func (s *MyTestSuite) SetupTest() {
	log.Println("SetupTest()")

	// Reset mock，避免在不同測試方法之間互相影響
	mockDatabaseRepository := repository.NewMockDatabaseRepository(s.T())
	s.examService.databaseRepository = mockDatabaseRepository
	s.mockDatabaseRepository = mockDatabaseRepository
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
	// Setup
	s.mockDatabaseRepository.EXPECT().
		CreateExam(mock.Anything, mock.Anything).
		Return("exam01", nil)

	// Test
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
	// Setup
	userId := "user01"
	id := primitive.NewObjectID()
	examId := id.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetExamById(mock.Anything, examId).
		Return(&model.Exam{
			Id:          id,
			Topic:       "topic01",
			Description: "desc01",
			Tags:        []string{"t01"},
			IsPublic:    true,
			UserId:      userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		UpdateExam(mock.Anything, mock.Anything).
		Return(nil)

	// Test
	resultExamId, err := s.examService.UpdateExam(
		examId,
		"topic02",
		"desc02",
		false,
		userId,
	)
	s.Nil(err)
	s.Equal(examId, resultExamId)
}

func (s *MyTestSuite) TestFindExams() {
	// Setup
	userId := "user01"

	var pageIndex int32 = 0
	var pageSize int32 = 10
	s.mockDatabaseRepository.EXPECT().
		FindExamsByUserIdOrderByUpdateAtDesc(
			mock.Anything,
			userId,
			pageIndex*pageSize,
			pageSize,
		).
		Return([]model.Exam{
			{},
			{},
			{},
		}, nil)

	var expectedTotal int32 = 3
	s.mockDatabaseRepository.EXPECT().
		CountExamsByUserId(mock.Anything, userId).
		Return(expectedTotal, nil)

	// Test
	total, pageCount, exams, err := s.examService.FindExams(pageIndex, pageSize, userId)
	s.Nil(err)
	s.Equal(expectedTotal, total)
	s.Equal(int32(1), pageCount)
	s.Equal(3, len(exams))
}

func (s *MyTestSuite) TestDeleteExam() {
	// Setup
	userId := "user01"

	id := primitive.NewObjectID()
	examId := id.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetExamById(mock.Anything, examId).
		Return(&model.Exam{
			Id:          id,
			Topic:       "topic01",
			Description: "desc01",
			Tags:        []string{"t01"},
			IsPublic:    true,
			UserId:      userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		WithTransaction(mock.Anything).
		Return(nil, nil)

	// Test
	err := s.examService.DeleteExam(examId, userId)
	s.Nil(err)
}

func (s *MyTestSuite) TestCreateQuestion() {
	// Setup
	userId := "user01"
	expectedQuestionId := "q01"
	id := primitive.NewObjectID()
	examId := id.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetExamById(mock.Anything, examId).
		Return(&model.Exam{
			Id:          id,
			Topic:       "topic01",
			Description: "desc01",
			Tags:        []string{"t01"},
			IsPublic:    true,
			UserId:      userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		CreateQuestion(mock.Anything, mock.Anything).
		Return(expectedQuestionId, nil)

	// Test
	questionId, err := s.examService.CreateQuestion(
		id.Hex(),
		"ask",
		[]string{"a01", "a02"},
		userId)
	s.Nil(err)
	s.Equal(expectedQuestionId, questionId)
}

func (s *MyTestSuite) TestUpdaetQuestion() {
	// Setup
	userId := "user01"

	id := primitive.NewObjectID()
	questionId := id.Hex()
	examId := "exam01"

	s.mockDatabaseRepository.EXPECT().
		GetQuestionById(mock.Anything, questionId).
		Return(&model.Question{
			Id:      id,
			ExamId:  examId,
			Ask:     "ask01",
			Answers: []string{"a01", "a02"},
			UserId:  userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		WithTransaction(mock.Anything).
		Return(nil, nil)

	// Test
	resultQuestionId, err := s.examService.UpdateQuestion(
		questionId,
		"ask02",
		[]string{"b01", "b02"},
		userId)
	s.Nil(err)
	s.Equal(questionId, resultQuestionId)
}

func (s *MyTestSuite) TestFindQuestions() {
	// Setup
	userId := "user01"
	var pageIndex int32 = 2
	var pageSize int32 = 10
	skip := pageIndex * pageSize
	limit := pageSize
	var expectedTotal int32 = 23

	id := primitive.NewObjectID()
	examId := id.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetExamById(mock.Anything, examId).
		Return(&model.Exam{
			Id:          id,
			Topic:       "topic01",
			Description: "desc01",
			Tags:        []string{"t01"},
			IsPublic:    true,
			UserId:      userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		FindQuestionsByExamIdOrderByUpdateAtDesc(mock.Anything, examId, skip, limit).
		Return([]model.Question{
			{},
			{},
			{},
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		CountQuestionsByExamId(mock.Anything, examId).
		Return(expectedTotal, nil)

	// Test
	total, pageCount, questions, err := s.examService.FindQuestions(
		pageIndex,
		pageSize,
		examId,
		userId,
	)
	s.Nil(err)
	s.Equal(expectedTotal, total)
	s.Equal(int32(3), pageCount)
	s.Equal(3, len(questions))
}

func (s *MyTestSuite) TestDeleteQuestion() {
	// Setup
	userId := "user01"
	id := primitive.NewObjectID()
	questionId := id.Hex()
	examId := "exam01"

	s.mockDatabaseRepository.EXPECT().
		GetQuestionById(mock.Anything, questionId).
		Return(&model.Question{
			Id:      id,
			ExamId:  examId,
			Ask:     "ask02",
			Answers: []string{"b01", "b02"},
			UserId:  userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		WithTransaction(mock.Anything).
		Return(nil, nil)

	// Test
	err := s.examService.DeleteQuestion(questionId, userId)
	s.Nil(err)
}

func (s *MyTestSuite) TestFindRandomQuestions() {
	// Setup
	id := primitive.NewObjectID()
	examId := id.Hex()
	userId := "user01"
	var size int32 = 10
	var limit int32 = 1
	var questionCount int32 = 5

	s.mockDatabaseRepository.EXPECT().
		GetExamById(mock.Anything, examId).
		Return(&model.Exam{
			Id:          id,
			Topic:       "topic01",
			Description: "desc01",
			Tags:        []string{"t01"},
			IsPublic:    true,
			UserId:      userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		FindQuestionsByExamIdOrderByUpdateAtDesc(mock.Anything, examId, mock.Anything, limit).
		Return([]model.Question{
			{},
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		CountQuestionsByExamId(mock.Anything, examId).
		Return(questionCount, nil)

	// Test
	exam, questions, err := s.examService.FindRandomQuestions(
		examId,
		userId,
		size,
	)
	s.Nil(err)
	s.NotEmpty(exam)
	s.EqualValues(questionCount, len(questions))
}

func (s *MyTestSuite) TestCreateExamRecord() {
	// Setup
	userId := "user01"
	id := primitive.NewObjectID()
	examId := id.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetExamById(mock.Anything, examId).
		Return(&model.Exam{
			Id:          id,
			Topic:       "topic",
			Description: "desc",
			Tags:        []string{"a01"},
			IsPublic:    true,
			UserId:      userId,
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		WithTransaction(mock.Anything).
		Return(nil, nil)

	// Test
	err := s.examService.CreateExamRecord(
		examId,
		10,
		[]string{"question01", "question02", "quesiton03"},
		userId,
	)
	s.Nil(err)
}

func (s *MyTestSuite) TestFindExamInfos() {
	// Setup
	userId := "user01"
	isPublic := true

	examId1 := "658875894c61d5f50a71a7b6"
	id1, err := primitive.ObjectIDFromHex(examId1)
	s.Nil(err)

	examId2 := "658875a9512667185df5e0b9"
	id2, err := primitive.ObjectIDFromHex(examId2)
	s.Nil(err)

	var (
		questionCount1 int32 = 10
		questionCount2 int32 = 20
		recordCount1   int32 = 1
		recordCount2   int32 = 2
	)

	s.mockDatabaseRepository.EXPECT().
		FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(mock.Anything, userId, isPublic).
		Return([]model.Exam{
			{
				Id:          id1,
				Topic:       "topic01",
				Description: "desc01",
				Tags:        []string{"a01"},
				IsPublic:    isPublic,
				UserId:      userId,
			},
			{
				Id:          id2,
				Topic:       "topic02",
				Description: "desc02",
				Tags:        []string{"a02"},
				IsPublic:    isPublic,
				UserId:      userId,
			},
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		CountQuestionsByExamId(mock.Anything, examId1).
		Return(questionCount1, nil)
	s.mockDatabaseRepository.EXPECT().
		CountQuestionsByExamId(mock.Anything, examId2).
		Return(questionCount2, nil)
	s.mockDatabaseRepository.EXPECT().
		CountExamRecordsByExamIdAndUserId(mock.Anything, examId1, userId).
		Return(recordCount1, nil)
	s.mockDatabaseRepository.EXPECT().
		CountExamRecordsByExamIdAndUserId(mock.Anything, examId2, userId).
		Return(recordCount2, nil)

	// Test
	examInfos, err := s.examService.FindExamInfos(userId, isPublic)
	s.Nil(err)
	s.Equal(2, len(examInfos))
	s.Equal(questionCount1, examInfos[0].QuestionCount)
	s.Equal(questionCount2, examInfos[1].QuestionCount)
	s.Equal(recordCount1, examInfos[0].RecordCount)
	s.Equal(recordCount2, examInfos[1].RecordCount)
}
