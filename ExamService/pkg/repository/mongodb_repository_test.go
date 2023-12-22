package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

type MyTestSuite struct {
	suite.Suite
	repo                    DatabaseRepository
	uri                     string
	ctx                     context.Context
	mongodbContainer        *mongodb.MongoDBContainer
	examIdForTestGetExam    string
	examIdForTestUpdateExam string
	examIdForTestDeleteExam string
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

	// 建立測試資料
	err = createTestData(s, uri)
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

	// 不呼叫 panic，為了繼續往下執行去關閉 container
	if err := s.repo.DisconnectDB(); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

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
		UserId:      "user01",
	})
	s.Nil(err)
	s.NotEmpty(examId)
}

func (s *MyTestSuite) TestUpdateExam() {
	id, err := primitive.ObjectIDFromHex(s.examIdForTestUpdateExam)
	s.Nil(err)

	ctx := context.TODO()
	err = s.repo.UpdateExam(ctx, model.Exam{
		Id:          id,
		Topic:       "TestUpdateExam_u01",
		Description: "desc01_u01",
		IsPublic:    false,
		Tags:        []string{"aaa", "bbb", "ccc"},
		UserId:      "user01",
	})
	s.Nil(err)
}

func (s *MyTestSuite) TestGetExam() {
	ctx := context.TODO()
	exam, err := s.repo.GetExam(ctx, s.examIdForTestGetExam)
	s.Nil(err)
	s.NotNil(exam)
}

func (s *MyTestSuite) TestFindExamsOrderByUpdateAtDesc() {
	ctx := context.TODO()
	exams, err := s.repo.FindExamsOrderByUpdateAtDesc(ctx, "user01", 10, 10)
	s.Nil(err)
	s.NotEmpty(exams)
	s.Equal(10, len(exams))
}

func (s *MyTestSuite) TestDeleteExam() {
	ctx := context.TODO()
	err := s.repo.DeleteExam(ctx, s.examIdForTestDeleteExam)
	s.Nil(err)
}

func createTestData(s *MyTestSuite, uri string) error {
	ctx := context.TODO()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		return err
	}

	exams := []interface{}{}

	for i := 0; i < 30; i++ {
		exams = append(exams, model.Exam{
			Topic:       fmt.Sprintf("topic_%d", i),
			Description: "jsut for test",
			Tags:        []string{"tag01", "tag02"},
			IsPublic:    true,
			UserId:      "user01",
		})
	}
	collection := client.Database("learnEnglish").Collection("exams")
	_, err = collection.InsertMany(ctx, exams)
	if err != nil {
		return err
	}

	// For TestGetExam
	result, err := collection.InsertOne(ctx, model.Exam{
		Topic:       "TestGetExam",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      "user01",
	})
	if err != nil {
		return err
	}
	s.examIdForTestGetExam = result.InsertedID.(primitive.ObjectID).Hex()

	// For TestUpdateExam
	result, err = collection.InsertOne(ctx, model.Exam{
		Topic:       "TestUpdateExam",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      "user01",
	})
	if err != nil {
		return err
	}
	s.examIdForTestUpdateExam = result.InsertedID.(primitive.ObjectID).Hex()

	// For TestDeleteExam
	result, err = collection.InsertOne(ctx, model.Exam{
		Topic:       "TestDeleteExam",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      "user01",
	})
	if err != nil {
		return err
	}
	s.examIdForTestDeleteExam = result.InsertedID.(primitive.ObjectID).Hex()

	if err := client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
