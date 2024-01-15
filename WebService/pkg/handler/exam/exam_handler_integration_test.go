package exam

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/examservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

type MyIntegrationTestSuite struct {
	suite.Suite
	examHandler           examHandler
	compose               tc.ComposeStack
	client                *mongo.Client
	userCollection        *mongo.Collection
	examCollection        *mongo.Collection
	questionCollection    *mongo.Collection
	answerWrongCollection *mongo.Collection
	examRecordCollection  *mongo.Collection
	userId                string
	adminUserId           string
}

func TestMyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MyIntegrationTestSuite))
}

// run once, before test suite methods
func (s *MyIntegrationTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	// Setup ExamService
	compose, err := tc.NewDockerCompose("./docker-compose.yml")
	if err != nil {
		s.FailNow(err.Error())
	}

	s.compose = compose

	ctx, cancel := context.WithCancel(context.Background())
	s.T().Cleanup(cancel)
	err = compose.
		WaitForService("exam-service", wait.NewLogStrategy("Starting gRPC server at")).
		Up(ctx, tc.Wait(true))
	if err != nil {
		s.FailNow(err.Error())
	}

	databaseName := "Test_LearnEnglish"
	mongoDBURI := "mongodb://127.0.0.1:27017"

	// 連線到資料庫
	databaseRepository := repository.NewMongoDBRepository(databaseName)
	err = databaseRepository.ConnectDB(ctx, mongoDBURI)
	if err != nil {
		s.FailNow(err.Error())
	}

	// ExamService
	examService := examservice.New(":8090")
	err = examService.Connect()
	if err != nil {
		s.FailNow(err.Error())
	}

	// 用來建立測試資料的 client
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(mongoDBURI).SetTimeout(10*time.Second),
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.client = client
	s.userCollection = client.Database(databaseName).Collection("users")
	s.examCollection = client.Database(databaseName).Collection("exams")
	s.questionCollection = client.Database(databaseName).Collection("questions")
	s.answerWrongCollection = client.Database(databaseName).Collection("answerwrongs")
	s.examRecordCollection = client.Database(databaseName).Collection("examrecords")

	// 新增測試資料
	now := time.Now()
	result, err := s.userCollection.InsertOne(ctx, model.User{
		Username:  "test-admin",
		Password:  "test123",
		Role:      "admin",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		s.FailNow(err.Error())
	}

	s.adminUserId = result.InsertedID.(primitive.ObjectID).Hex()

	username := "user01"
	role := "user"
	result, err = s.userCollection.InsertOne(ctx, model.User{
		Username:  username,
		Password:  "test123",
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		s.FailNow(err.Error())
	}

	s.userId = result.InsertedID.(primitive.ObjectID).Hex()

	// Mock JWT
	utilGetJWTClaims = func(c echo.Context) *util.JwtCustomClaims {
		return &util.JwtCustomClaims{
			UserId:           s.userId,
			Username:         username,
			Role:             role,
			RegisteredClaims: jwt.RegisteredClaims{},
		}
	}

	s.examHandler = examHandler{
		examService:        examService,
		databaseRepository: databaseRepository,
	}
}

// run once, after test suite methods
func (s *MyIntegrationTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
	ctx := context.Background()

	if err := s.client.Disconnect(ctx); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	if err := s.examHandler.databaseRepository.DisconnectDB(ctx); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// 程式結束時，結束微服務連線
	if err := s.examHandler.examService.Disconnect(); err != nil {
		log.Printf("examService disconnect() error: %v", err)
	}

	// 終止 container
	if err := s.compose.Down(ctx, tc.RemoveOrphans(true), tc.RemoveImagesLocal); err != nil {
		log.Printf("compose.Down() error: %v", err)
	}
}

// run before each test
func (s *MyIntegrationTestSuite) SetupTest() {
	log.Println("SetupTest()")
}

// run after each test
func (s *MyIntegrationTestSuite) TearDownTest() {
	log.Println("TearDownTest()")
}

// run before each test
func (s *MyIntegrationTestSuite) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
}

// run after each test
func (s *MyIntegrationTestSuite) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}

func (s *MyIntegrationTestSuite) TestCreateExam() {
	// Setup
	requestJSON := `{
  	"topic": "t01",
  	"description": "d01",
  	"isPublic": false
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err := s.examHandler.CreateExam(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindExams() {
	// Setup
	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err := s.examHandler.FindExams(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestUpdateExam() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	requestJSON := `{
		"_id": "` + examId + `",
  	"topic": "t01",
  	"description": "d01",
  	"isPublic": false
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err = s.examHandler.UpdateExam(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindQuestions() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Question{
			ExamId:    examId,
			Ask:       fmt.Sprintf("Question_%d", i),
			Answers:   []string{"a01", "a02"},
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	_, err = s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	// Test
	err = s.examHandler.FindQuestions(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestCreateQuestion() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	requestJSON := `{
  	"ask": "ask01",
  	"answers": ["a01", "a02"]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	// Test
	err = s.examHandler.CreateQuestion(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindRandomQuestions() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Question{
			ExamId:    examId,
			Ask:       fmt.Sprintf("Question_%d", i),
			Answers:   []string{"a01", "a02"},
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	_, err = s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	// Test
	err = s.examHandler.FindRandomQuestions(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindExamRecordOverview() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Question{
			ExamId:    examId,
			Ask:       fmt.Sprintf("Question_%d", i),
			Answers:   []string{"a01", "a02"},
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	questionResult, err := s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	documents = []interface{}{}

	for i := 0; i < len(questionResult.InsertedIDs); i++ {
		questionId := questionResult.InsertedIDs[i].(primitive.ObjectID).Hex()
		documents = append(documents, model.AnswerWrong{
			ExamId:     examId,
			QuestionId: questionId,
			Times:      3,
			UserId:     s.userId,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	documents = []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId:    examId,
			Score:     10,
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	_, err = s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	// Test
	err = s.examHandler.FindExamRecordOverview(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindExamRecords() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Question{
			ExamId:    examId,
			Ask:       fmt.Sprintf("Question_%d", i),
			Answers:   []string{"a01", "a02"},
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	questionResult, err := s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	documents = []interface{}{}

	for i := 0; i < len(questionResult.InsertedIDs); i++ {
		questionId := questionResult.InsertedIDs[i].(primitive.ObjectID).Hex()
		documents = append(documents, model.AnswerWrong{
			ExamId:     examId,
			QuestionId: questionId,
			Times:      3,
			UserId:     s.userId,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	documents = []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId:    examId,
			Score:     10,
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	_, err = s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	// Test
	err = s.examHandler.FindExamRecords(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindExamInfosWhenNotSignIn() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.adminUserId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Question{
			ExamId:    examId,
			Ask:       fmt.Sprintf("Question_%d", i),
			Answers:   []string{"a01", "a02"},
			UserId:    s.adminUserId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	questionResult, err := s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	documents = []interface{}{}

	for i := 0; i < len(questionResult.InsertedIDs); i++ {
		questionId := questionResult.InsertedIDs[i].(primitive.ObjectID).Hex()
		documents = append(documents, model.AnswerWrong{
			ExamId:     examId,
			QuestionId: questionId,
			Times:      3,
			UserId:     s.adminUserId,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	documents = []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId:    examId,
			Score:     10,
			UserId:    s.adminUserId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	_, err = s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err = s.examHandler.FindExamInfosWhenNotSignIn(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindExamInfosWhenSignIn() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.examCollection.InsertOne(ctx, model.Exam{
		Topic:       "topic01",
		Description: "jsut for test",
		Tags:        []string{"tag01", "tag02"},
		IsPublic:    true,
		UserId:      s.userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	s.Nil(err)

	examId := result.InsertedID.(primitive.ObjectID).Hex()

	documents := []interface{}{}
	size := 10

	for i := 0; i < size; i++ {
		documents = append(documents, model.Question{
			ExamId:    examId,
			Ask:       fmt.Sprintf("Question_%d", i),
			Answers:   []string{"a01", "a02"},
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	questionResult, err := s.questionCollection.InsertMany(ctx, documents)
	s.Nil(err)

	documents = []interface{}{}

	for i := 0; i < len(questionResult.InsertedIDs); i++ {
		questionId := questionResult.InsertedIDs[i].(primitive.ObjectID).Hex()
		documents = append(documents, model.AnswerWrong{
			ExamId:     examId,
			QuestionId: questionId,
			Times:      3,
			UserId:     s.userId,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	documents = []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.ExamRecord{
			ExamId:    examId,
			Score:     10,
			UserId:    s.userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	_, err = s.examRecordCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err = s.examHandler.FindExamInfosWhenSignIn(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}
