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
	repo                  DatabaseRepository
	uri                   string
	ctx                   context.Context
	mongodbContainer      *mongodb.MongoDBContainer
	client                *mongo.Client
	examCollection        *mongo.Collection
	questionCollection    *mongo.Collection
	answerWrongCollection *mongo.Collection
	examRecordCollection  *mongo.Collection
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

	s.repo = NewMongoDBRepository(DATABASE)
	s.repo.ConnectDB(uri)
	s.uri = uri
	s.ctx = ctx
	s.mongodbContainer = mongodbContainer

	// 用來建立測試資料的 client
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		panic(err)
	}

	s.client = client
	s.examCollection = client.Database(DATABASE).Collection("exams")
	s.questionCollection = client.Database(DATABASE).Collection("questions")
	s.answerWrongCollection = s.client.Database(DATABASE).Collection("answerwrongs")
	s.examRecordCollection = s.client.Database(DATABASE).Collection("examrecords")
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")

	// 不呼叫 panic，為了繼續往下執行去關閉 container
	if err := s.client.Disconnect(context.TODO()); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

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
	ctx := context.TODO()

	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "TestUpdateExam",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      "user01",
	})
	s.Nil(err)
	id := result.InsertedID.(primitive.ObjectID)

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

func (s *MyTestSuite) TestGetExamById() {
	ctx := context.TODO()

	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "TestGetExamById",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      "user01",
	})
	s.Nil(err)
	examId := result.InsertedID.(primitive.ObjectID).Hex()

	exam, err := s.repo.GetExamById(ctx, examId)
	s.Nil(err)
	s.NotNil(exam)
}

func (s *MyTestSuite) TestFindExamsByUserIdOrderByUpdateAtDesc() {
	ctx := context.TODO()

	userId := "user01"
	documents := []interface{}{}

	for i := 0; i < 30; i++ {
		documents = append(documents, model.Exam{
			Topic:       fmt.Sprintf("topic_%d", i),
			Description: "jsut for test",
			Tags:        []string{"tag01", "tag02"},
			IsPublic:    true,
			UserId:      userId,
		})
	}
	_, err := s.examCollection.InsertMany(ctx, documents)
	s.Nil(err)

	exams, err := s.repo.FindExamsByUserIdOrderByUpdateAtDesc(ctx, userId, 10, 10)
	s.Nil(err)
	s.NotEmpty(exams)
	s.Equal(10, len(exams))
}

func (s *MyTestSuite) TestFindExamsByUserIdAndIsPublicOrderByUpdateAtDesc() {
	ctx := context.TODO()

	userId := "TestFindExamsByUserIdAndIsPublicOrderByUpdateAtDesc"
	isPublic := true
	size := 10
	documents := []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.Exam{
			Topic:       fmt.Sprintf("topic_%d", i),
			Description: "jsut for test",
			Tags:        []string{"tag01", "tag02"},
			IsPublic:    isPublic,
			UserId:      userId,
		})
	}
	_, err := s.examCollection.InsertMany(ctx, documents)
	s.Nil(err)

	exams, err := s.repo.FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(ctx, userId, isPublic)
	s.Nil(err)
	s.Equal(size, len(exams))
}

func (s *MyTestSuite) TestDeleteExamById() {
	ctx := context.TODO()

	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "TestDeleteExamById",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      "user01",
	})
	s.Nil(err)
	examId := result.InsertedID.(primitive.ObjectID).Hex()

	deletedCount, err := s.repo.DeleteExamById(ctx, examId)
	s.Nil(err)
	s.Equal(int64(1), deletedCount)
}

func (s *MyTestSuite) TestCountExamsByUserId() {
	ctx := context.TODO()

	userId := "user01"
	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Exam{
			Topic:       fmt.Sprintf("topic_%d", i),
			Description: "jsut for test",
			Tags:        []string{"tag01", "tag02"},
			IsPublic:    true,
			UserId:      userId,
		})
	}
	_, err := s.examCollection.InsertMany(ctx, documents)
	s.Nil(err)

	count, err := s.repo.CountExamsByUserId(ctx, userId)
	s.Nil(err)
	s.Equal(int64(size), count)
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
	ctx := context.TODO()

	result, err := s.questionCollection.InsertOne(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "TestUpdateQuestion",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	s.Nil(err)
	id := result.InsertedID.(primitive.ObjectID)

	err = s.repo.UpdateQuestion(ctx, model.Question{
		Id:      id,
		ExamId:  "exam02",
		Ask:     "TestUpdateQuestion02",
		Answers: []string{"a011", "a022"},
		UserId:  "user01",
	})
	s.Nil(err)
}

func (s *MyTestSuite) TestGetQuestionById() {
	ctx := context.TODO()

	result, err := s.questionCollection.InsertOne(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "TestGetQuestion",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	s.Nil(err)
	questionId := result.InsertedID.(primitive.ObjectID).Hex()

	question, err := s.repo.GetQuestionById(ctx, questionId)
	s.Nil(err)
	s.NotNil(question)
}

func (s *MyTestSuite) TestFindQuestionsByExamIdAndUserIdOrderByUpdateAtDesc() {
	ctx := context.TODO()

	examId := "TestFindQuestionsByExamIdAndUserIdOrderByUpdateAtDesc"
	userId := "user01"
	documents := []interface{}{}

	for i := 0; i < 30; i++ {
		documents = append(documents, model.Question{
			ExamId:  examId,
			Ask:     fmt.Sprintf("Question_%d", i),
			Answers: []string{"a01", "a02"},
			UserId:  userId,
		})
	}
	_, err := s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	questions, err := s.repo.FindQuestionsByExamIdAndUserIdOrderByUpdateAtDesc(
		ctx,
		examId,
		userId,
		10,
		10,
	)
	s.Nil(err)
	s.NotEmpty(questions)
	s.Equal(10, len(questions))
}

func (s *MyTestSuite) TestDeleteQuestionById() {
	ctx := context.TODO()

	result, err := s.questionCollection.InsertOne(ctx, model.Question{
		ExamId:  "exam01",
		Ask:     "TestDeleteQuestion",
		Answers: []string{"a01", "a02"},
		UserId:  "user01",
	})
	s.Nil(err)
	questionId := result.InsertedID.(primitive.ObjectID).Hex()

	deletedCount, err := s.repo.DeleteQuestionById(ctx, questionId)
	s.Nil(err)
	s.Equal(int64(1), deletedCount)
}

func (s *MyTestSuite) TestDeleteQuestionsByExamId() {
	ctx := context.TODO()

	examId := "TestDeleteQuestionsByExamId"
	size := 10
	questions := []interface{}{}

	for i := 0; i < size; i++ {
		questions = append(questions, model.Question{
			ExamId:  examId,
			Ask:     fmt.Sprintf("Question_%d", i),
			Answers: []string{"a01", "a02"},
			UserId:  "user01",
		})
	}
	_, err := s.questionCollection.InsertMany(ctx, questions)
	s.Nil(err)

	deletedCount, err := s.repo.DeleteQuestionsByExamId(ctx, examId)
	s.Nil(err)
	s.Equal(int64(size), deletedCount)
}

func (s *MyTestSuite) TestCountQuestionsByExamIdAndUserId() {
	ctx := context.TODO()

	examId := "TestCountQuestionsByExamIdAndUserId"
	userId := "user01"
	size := 10
	questions := []interface{}{}

	for i := 0; i < size; i++ {
		questions = append(questions, model.Question{
			ExamId:  examId,
			Ask:     fmt.Sprintf("Question_%d", i),
			Answers: []string{"a01", "a02"},
			UserId:  userId,
		})
	}
	_, err := s.questionCollection.InsertMany(ctx, questions)
	s.Nil(err)

	count, err := s.repo.CountQuestionsByExamIdAndUserId(ctx, examId, userId)
	s.Nil(err)
	s.Equal(int64(size), count)
}

func (s *MyTestSuite) TestDeleteAnswerWrongsByQuestionId() {
	ctx := context.TODO()

	questionId := "TestDeleteAnswerWrongsByQuestionId"
	size := 10
	documents := []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.AnswerWrong{
			ExamId:     "exam_abc_01",
			QuestionId: questionId,
			Times:      10,
			UserId:     "user01",
		})
	}

	_, err := s.answerWrongCollection.InsertMany(ctx, documents)
	s.Nil(err)

	deletedCount, err := s.repo.DeleteAnswerWrongsByQuestionId(ctx, questionId)
	s.Nil(err)
	s.Equal(int64(size), deletedCount)
}

func (s *MyTestSuite) TestDeleteAnswerWrongsByExamId() {
	ctx := context.TODO()

	examId := "TestDeleteAnswerWrongsByExamId"
	size := 10
	documents := []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.AnswerWrong{
			ExamId:     examId,
			QuestionId: "q01",
			Times:      10,
			UserId:     "user01",
		})
	}
	_, err := s.answerWrongCollection.InsertMany(ctx, documents)
	s.Nil(err)

	deletedCount, err := s.repo.DeleteAnswerWrongsByExamId(ctx, examId)
	s.Nil(err)
	s.Equal(int64(size), deletedCount)
}

func (s *MyTestSuite) TestUpsertAnswerWrongByTimesPlusOne() {
	ctx := context.TODO()

	examId := "TestUpsertAnswerWrongByTimesPlusOne"
	questionId := "TestUpsertAnswerWrongByTimesPlusOne_q01"
	userId := "TestUpsertAnswerWrongByTimesPlusOne_u01"

	modifiedCount, upsertedCount, err := s.repo.UpsertAnswerWrongByTimesPlusOne(
		ctx,
		examId,
		questionId,
		userId,
	)
	s.Nil(err)
	s.Equal(int64(0), modifiedCount)
	s.Equal(int64(1), upsertedCount)

	modifiedCount, upsertedCount, err = s.repo.UpsertAnswerWrongByTimesPlusOne(
		ctx,
		examId,
		questionId,
		userId,
	)
	s.Nil(err)
	s.Equal(int64(1), modifiedCount)
	s.Equal(int64(0), upsertedCount)
}

func (s *MyTestSuite) TestDeleteExamRecordsByExamId() {
	ctx := context.TODO()

	examId := "TestDeleteExamRecordsByExamId"
	size := 10
	documents := []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId: examId,
			Score:  6,
			UserId: "user01",
		})
	}
	_, err := s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	deletedCount, err := s.repo.DeleteExamRecordsByExamId(ctx, examId)
	s.Nil(err)
	s.Equal(int64(size), deletedCount)
}

func (s *MyTestSuite) TestCreateExamRecord() {
	ctx := context.TODO()

	examRecordId, err := s.repo.CreateExamRecord(ctx, model.ExamRecord{
		ExamId: "TestCreateExamRecord",
		Score:  10,
		UserId: "user01",
	})
	s.Nil(err)
	s.NotEmpty(examRecordId)
}

func (s *MyTestSuite) TestFindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc() {
	ctx := context.TODO()

	examId := "exam01"
	userId := "user01"
	documents := []interface{}{}

	for i := 0; i < 30; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId: examId,
			Score:  10,
			UserId: userId,
		})
	}
	_, err := s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	examRecords, err := s.repo.FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
		ctx,
		examId,
		userId,
		10,
		10,
	)
	s.Nil(err)
	s.Equal(10, len(examRecords))
}

func (s *MyTestSuite) TestCountExamRecordsByExamIdAndUserId() {
	ctx := context.TODO()

	examId := "exam01"
	userId := "user01"
	size := 10
	documents := []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId: examId,
			Score:  10,
			UserId: userId,
		})
	}
	_, err := s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	count, err := s.repo.CountExamRecordsByExamIdAndUserId(ctx, examId, userId)
	s.Nil(err)
	s.Equal(int64(size), count)
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
