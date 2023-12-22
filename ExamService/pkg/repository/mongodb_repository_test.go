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

// 使用測試的資料庫
const DATABASE = "learnEnglish_test"

type MyTestSuite struct {
	suite.Suite
	repo             DatabaseRepository
	uri              string
	ctx              context.Context
	mongodbContainer *mongodb.MongoDBContainer

	examIdForTestUpdateExam string
	examIdForTestGetExam    string
	examIdForTestDeleteExam string

	questionIdForTestUpdateQuestion      string
	questionIdForTestGetQuestion         string
	questionIdForTestDeleteQuestion      string
	examIdForTestDeleteQuestionsByExamId string

	questionIdForTestDeleteAnswerWrongByQuestionId string
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

	s.repo = NewMongoDBRepository(DATABASE)
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
	repo := NewMongoDBRepository(DATABASE)

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

func (s *MyTestSuite) TestCreateQuestion() {
	ctx := context.TODO()
	questionId, err := s.repo.CreateQuestion(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "ask01",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	s.Nil(err)
	s.NotEmpty(questionId)
}

func (s *MyTestSuite) TestUpdateQuestion() {
	id, err := primitive.ObjectIDFromHex(s.questionIdForTestUpdateQuestion)
	s.Nil(err)

	ctx := context.TODO()
	err = s.repo.UpdateQuestion(ctx, model.Question{
		Id:      id,
		ExamId:  "examId",
		Ask:     "TestUpdateQuestion_u01",
		Answers: []string{"a011", "a022"},
		UserId:  "user01",
	})
	s.Nil(err)
}

func (s *MyTestSuite) TestGetQuestion() {
	ctx := context.TODO()
	question, err := s.repo.GetQuestion(ctx, s.questionIdForTestGetQuestion)
	s.Nil(err)
	s.NotNil(question)
}

func (s *MyTestSuite) TestFindQuestionsOrderByUpdateAtDesc() {
	ctx := context.TODO()
	questions, err := s.repo.FindQuestionsOrderByUpdateAtDesc(ctx, "exam01", 10, 10)
	s.Nil(err)
	s.NotEmpty(questions)
	s.Equal(10, len(questions))
}

func (s *MyTestSuite) TestDeleteQuestion() {
	ctx := context.TODO()
	err := s.repo.DeleteQuestion(ctx, s.questionIdForTestDeleteQuestion)
	s.Nil(err)
}

func (s *MyTestSuite) TestDeleteQuestionsByExamId() {
	ctx := context.TODO()
	err := s.repo.DeleteQuestionsByExamId(ctx, s.examIdForTestDeleteQuestionsByExamId)
	s.Nil(err)
}

func (s *MyTestSuite) TestDeleteAnswerWrongByQuestionId() {
	ctx := context.TODO()
	err := s.repo.DeleteAnswerWrongByQuestionId(
		ctx,
		s.questionIdForTestDeleteAnswerWrongByQuestionId,
	)
	s.Nil(err)
}

// testcontainers mongodb 不支援交易功能，所以註解此測試
// func (s *MyTestSuite) TestWithTransaction() {
// 	userId := "user_mongodb_test_002"
//
// 	_, err := s.repo.WithTransaction(func(ctx context.Context) (interface{}, error) {
// 		deleteExamId := ""
//
// 		// 新增 10 筆資料
// 		for i := 0; i < 10; i++ {
// 			examId, err := s.repo.CreateExam(ctx, model.Exam{
// 				Topic:       fmt.Sprintf("TestWithTransaction_%d", i),
// 				Description: "jsut for test",
// 				Tags:        []string{"tag01", "tag02"},
// 				IsPublic:    true,
// 				UserId:      userId,
// 			})
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			if i == 0 {
// 				deleteExamId = examId
// 			}
// 		}
//
// 		// 刪除 1 筆資料
// 		err := s.repo.DeleteExam(ctx, deleteExamId)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		return nil, nil
// 	})
// 	s.Nil(err)
//
// 	// 查詢看看交易是否真的成功新增資料
// 	ctx := context.TODO()
// 	exams, err := s.repo.FindExamsOrderByUpdateAtDesc(ctx, userId, 0, 100)
// 	s.Nil(err)
// 	s.Equal(9, len(exams))
// }

func createTestData(s *MyTestSuite, uri string) error {
	ctx := context.TODO()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		return err
	}

	// Exam
	collection := client.Database(DATABASE).Collection("exams")

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
	_, err = collection.InsertMany(ctx, exams)
	if err != nil {
		return err
	}

	// For TestUpdateExam
	result, err := collection.InsertOne(ctx, model.Exam{
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

	// For TestGetExam
	result, err = collection.InsertOne(ctx, model.Exam{
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

	// Question
	questionCollection := client.Database(DATABASE).Collection("questions")

	questions := []interface{}{}

	for i := 0; i < 30; i++ {
		questions = append(questions, model.Question{
			ExamId:  "exam01",
			Ask:     fmt.Sprintf("Question_%d", i),
			Answers: []string{"a01", "a02"},
			UserId:  "user01",
		})
	}
	_, err = questionCollection.InsertMany(ctx, questions)
	if err != nil {
		return err
	}

	// For TestUpdateQuestion
	result, err = questionCollection.InsertOne(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "TestUpdateQuestion",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	if err != nil {
		return err
	}
	s.questionIdForTestUpdateQuestion = result.InsertedID.(primitive.ObjectID).Hex()

	// For TestGetQuestion
	result, err = questionCollection.InsertOne(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "TestGetQuestion",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	if err != nil {
		return err
	}
	s.questionIdForTestGetQuestion = result.InsertedID.(primitive.ObjectID).Hex()

	// For TestDeleteQuestion
	result, err = questionCollection.InsertOne(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "TestDeleteQuestion",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	if err != nil {
		return err
	}
	s.questionIdForTestDeleteQuestion = result.InsertedID.(primitive.ObjectID).Hex()

	// For TestDeleteQuestionsByExamId
	s.examIdForTestDeleteQuestionsByExamId = "TestDeleteQuestionsByExamId"
	questions = []interface{}{}

	for i := 0; i < 30; i++ {
		questions = append(questions, model.Question{
			ExamId:  s.examIdForTestDeleteQuestionsByExamId,
			Ask:     fmt.Sprintf("Question_%d", i),
			Answers: []string{"a01", "a02"},
			UserId:  "user01",
		})
	}
	_, err = questionCollection.InsertMany(ctx, questions)
	if err != nil {
		return err
	}

	// AnswerWrong
	answerWrongCollection := client.Database(DATABASE).Collection("answerwrongs")

	s.questionIdForTestDeleteAnswerWrongByQuestionId = "TestDeleteAnswerWrongsByQuestionId"

	_, err = answerWrongCollection.InsertOne(ctx, model.AnswerWrong{
		ExamId:     "exam_abc_01",
		QuestionId: s.questionIdForTestDeleteAnswerWrongByQuestionId,
		Times:      10,
		UserId:     "user01",
	})
	if err != nil {
		return err
	}

	if err := client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
