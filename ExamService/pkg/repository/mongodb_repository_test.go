package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type MyTestSuite struct {
	suite.Suite
	repo             DatabaseRepository
	uri              string
	ctx              context.Context
	mongodbContainer *mongodb.MongoDBContainer
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}

// run once, before test suite methods
func (s *MyTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	// Run container
	ctx := context.Background()
	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:6"))
	if err != nil {
		panic(err)
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	s.repo = NewMongoDBRepository()
	s.repo.ConnectDB(uri)
	s.uri = uri
	s.ctx = ctx
	s.mongodbContainer = mongodbContainer
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")

	// Terminate container
	if err := s.mongodbContainer.Terminate(s.ctx); err != nil {
		panic(err)
	}
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

func (s *MyTestSuite) TestConnectDBAndDisconnectDB() {
	repo := NewMongoDBRepository()

	err := repo.ConnectDB(s.uri)
	s.Nil(err)

	err = repo.DisconnectDB()
	s.Nil(err)
}

func (s *MyTestSuite) TestCreateExam() {
	ctx := context.TODO()

	examId, err := s.repo.CreateExam(ctx, model.Exam{
		Topic:       "TestCreateExam",
		Description: "desc01",
		IsPublic:    true,
		Tags:        []string{},
		UserId:      "ddr01",
	})
	s.Nil(err)
	s.NotEmpty(examId)
}

func (s *MyTestSuite) TestUpdateExam() {
	ctx := context.TODO()

	examId, err := s.repo.CreateExam(ctx, model.Exam{
		Topic:       "TestUpdateExam",
		Description: "desc01",
		IsPublic:    true,
		Tags:        []string{},
		UserId:      "ddr01",
	})
	s.Nil(err)
	s.NotEmpty(examId)

	id, err := primitive.ObjectIDFromHex(examId)
	s.Nil(err)

	err = s.repo.UpdateExam(ctx, model.Exam{
		Id:          id,
		Topic:       "TestUpdateExam_u01",
		Description: "desc01_u01",
		IsPublic:    false,
		Tags:        []string{"tag01", "tag02"},
		UserId:      "ddr01",
	})
	s.Nil(err)
}

func (s *MyTestSuite) TestGetExam() {
	ctx := context.TODO()

	examId, err := s.repo.CreateExam(ctx, model.Exam{
		Topic:       "TestGetExam",
		Description: "desc01",
		IsPublic:    true,
		Tags:        []string{},
		UserId:      "ddr01",
	})
	s.Nil(err)
	s.NotEmpty(examId)

	exam, err := s.repo.GetExam(ctx, examId)
	s.Nil(err)
	s.NotNil(exam)
}

func (s *MyTestSuite) TestFindExamsOrderByUpdateAtDesc() {
	ctx := context.TODO()

	for i := 0; i < 20; i++ {
		examId, err := s.repo.CreateExam(ctx, model.Exam{
			Topic:       fmt.Sprintf("TestFindExamsOrderByUpdateAtDesc_%d", i),
			Description: "desc01",
			IsPublic:    true,
			Tags:        []string{},
			UserId:      "ddr01",
		})
		s.Nil(err)
		s.NotEmpty(examId)
	}

	exams, err := s.repo.FindExamsOrderByUpdateAtDesc(ctx, "ddr01", 10, 10)
	s.Nil(err)
	s.NotEmpty(exams)
	s.Equal(10, len(exams))
}

func (s *MyTestSuite) TestDeleteExam() {
	ctx := context.TODO()

	examId, err := s.repo.CreateExam(ctx, model.Exam{
		Topic:       "TestDeleteExam",
		Description: "desc01",
		IsPublic:    true,
		Tags:        []string{},
		UserId:      "ddr01",
	})
	s.Nil(err)
	s.NotEmpty(examId)

	err = s.repo.DeleteExam(ctx, examId)
	s.Nil(err)
}
