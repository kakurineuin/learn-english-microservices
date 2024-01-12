package exam

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/examservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

const USER_ID = "user01"

type MyTestSuite struct {
	suite.Suite
	examHandler            examHandler
	mockExamService        *examservice.MockExamService
	mockDatabaseRepository *repository.MockDatabaseRepository
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}

// run once, before test suite methods
func (s *MyTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	// Mock JWT
	utilGetJWTClaims = func(c echo.Context) *util.JwtCustomClaims {
		return &util.JwtCustomClaims{
			UserId:           USER_ID,
			Username:         "test01",
			Role:             "user",
			RegisteredClaims: jwt.RegisteredClaims{},
		}
	}

	s.examHandler = examHandler{
		examService:        nil,
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
	mockExamService := examservice.NewMockExamService(s.T())
	s.examHandler.examService = mockExamService
	s.mockExamService = mockExamService

	mockDatabaseRepository := repository.NewMockDatabaseRepository(s.T())
	s.examHandler.databaseRepository = mockDatabaseRepository
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

	s.mockExamService.EXPECT().
		CreateExam("t01", "d01", false, USER_ID).
		Return(&pb.CreateExamResponse{
			ExamId: "exam01",
		}, nil)

	// Test
	err := s.examHandler.CreateExam(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"examId": "exam01"}`, rec.Body.String())
}

func (s *MyTestSuite) TestFindExams() {
	// Setup
	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockExamService.EXPECT().
		FindExams(int32(0), int32(10), USER_ID).
		Return(&pb.FindExamsResponse{
			Total:     1,
			PageCount: 1,
			Exams: []*pb.Exam{
				{
					Id:          "id01",
					Topic:       "t01",
					Description: "d01",
					Tags:        []string{},
					IsPublic:    true,
					UserId:      USER_ID,
					CreatedAt:   nil,
					UpdatedAt:   nil,
				},
			},
		}, nil)

	// Test
	err := s.examHandler.FindExams(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"total": 1, "pageCount": 1, "exams": [{
		"_id": "id01",
		"topic": "t01",
		"description": "d01",
		"tags": [],
		"isPublic": true,
		"userId": "`+USER_ID+`",
		"createdAt": null,
		"updatedAt": null
	}]}`, rec.Body.String())
}

func (s *MyTestSuite) TestUpdateExam() {
	// Setup
	requestJSON := `{
		"_id": "exam01",
  	"topic": "t01",
  	"description": "d01",
  	"isPublic": false
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockExamService.EXPECT().
		UpdateExam("exam01", "t01", "d01", false, USER_ID).
		Return(&pb.UpdateExamResponse{
			ExamId: "exam01",
		}, nil)

	// Test
	err := s.examHandler.UpdateExam(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"examId": "exam01"}`, rec.Body.String())
}

func (s *MyTestSuite) TestDeleteExam() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	s.mockExamService.EXPECT().
		DeleteExam(examId, USER_ID).
		Return(&pb.DeleteExamResponse{}, nil)

	// Test
	err := s.examHandler.DeleteExam(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Empty(rec.Body.String())
}

func (s *MyTestSuite) TestFindQuestions() {
	// Setup
	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	s.mockExamService.EXPECT().
		FindQuestions(int32(0), int32(10), examId, USER_ID).
		Return(&pb.FindQuestionsResponse{
			Total:     1,
			PageCount: 1,
			Questions: []*pb.Question{
				{
					Id:        "id01",
					ExamId:    examId,
					Ask:       "ask01",
					Answers:   []string{"a01", "a02"},
					UserId:    USER_ID,
					CreatedAt: nil,
					UpdatedAt: nil,
				},
			},
		}, nil)

	// Test
	err := s.examHandler.FindQuestions(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"total": 1, "pageCount": 1, "questions": [{
		"_id": "id01",
		"examId": "`+examId+`",
		"ask": "ask01",
		"answers": ["a01", "a02"],
		"userId": "`+USER_ID+`",
		"createdAt": null,
		"updatedAt": null
	}]}`, rec.Body.String())
}

func (s *MyTestSuite) TestCreateQuestion() {
	// Setup
	requestJSON := `{
  	"ask": "ask01",
  	"answers": ["a01", "a02"]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	s.mockExamService.EXPECT().
		CreateQuestion(examId, "ask01", []string{"a01", "a02"}, USER_ID).
		Return(&pb.CreateQuestionResponse{
			QuestionId: "question01",
		}, nil)

	// Test
	err := s.examHandler.CreateQuestion(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"questionId": "question01"}`, rec.Body.String())
}

func (s *MyTestSuite) TestUpdateQuestion() {
	// Setup
	requestJSON := `{
		"_id": "question01",
  	"ask": "ask01",
  	"answers": ["a01", "a02"]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockExamService.EXPECT().
		UpdateQuestion("question01", "ask01", []string{"a01", "a02"}, USER_ID).
		Return(&pb.UpdateQuestionResponse{
			QuestionId: "question01",
		}, nil)

	// Test
	err := s.examHandler.UpdateQuestion(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"questionId": "question01"}`, rec.Body.String())
}

func (s *MyTestSuite) TestDeleteQuestion() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	questionId := "question01"
	c.SetParamNames("questionId")
	c.SetParamValues(questionId)

	s.mockExamService.EXPECT().
		DeleteQuestion(questionId, USER_ID).
		Return(&pb.DeleteQuestionResponse{}, nil)

	// Test
	err := s.examHandler.DeleteQuestion(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Empty(rec.Body.String())
}

func (s *MyTestSuite) TestFindRandomQuestions() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)
	var size int32 = 10

	s.mockExamService.EXPECT().
		FindRandomQuestions(examId, USER_ID, size).
		Return(&pb.FindRandomQuestionsResponse{
			Exam: &pb.Exam{
				Id:          "examId01",
				Topic:       "t01",
				Description: "d01",
				Tags:        []string{},
				IsPublic:    true,
				UserId:      USER_ID,
				CreatedAt:   nil,
				UpdatedAt:   nil,
			},
			Questions: []*pb.Question{
				{
					Id:        "id01",
					ExamId:    examId,
					Ask:       "ask01",
					Answers:   []string{"a01", "a02"},
					UserId:    USER_ID,
					CreatedAt: nil,
					UpdatedAt: nil,
				},
			},
		}, nil)

	// Test
	err := s.examHandler.FindRandomQuestions(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"exam": {
		"_id": "examId01",
		"topic": "t01",
		"description": "d01",
		"tags": [],
		"isPublic": true,
		"userId": "`+USER_ID+`",
		"createdAt": null,
		"updatedAt": null
	}, 
	"questions": [{
		"_id": "id01",
		"examId": "`+examId+`",
		"ask": "ask01",
		"answers": ["a01", "a02"],
		"userId": "`+USER_ID+`",
		"createdAt": null,
		"updatedAt": null
	}]}`, rec.Body.String())
}

func (s *MyTestSuite) TestCreateExamRecord() {
	// Setup
	requestJSON := `{
  	"Score": 10,
  	"WrongQuestionIds": ["q01", "q02"]
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	s.mockExamService.EXPECT().
		CreateExamRecord(examId, int32(10), []string{"q01", "q02"}, USER_ID).
		Return(&pb.CreateExamRecordResponse{}, nil)

	// Test
	err := s.examHandler.CreateExamRecord(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{}`, rec.Body.String())
}

func (s *MyTestSuite) TestFindExamRecordOverview() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	s.mockExamService.EXPECT().
		FindExamRecordOverview(examId, USER_ID, mock.Anything).
		Return(&pb.FindExamRecordOverviewResponse{}, nil)

	// Test
	err := s.examHandler.FindExamRecordOverview(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(
		`{"startDate": "", "exam": null, "questions": [], "answerWrongs": [], "examRecords": []}`,
		rec.Body.String(),
	)
}

func (s *MyTestSuite) TestFindExamRecords() {
	// Setup
	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	examId := "exam01"
	c.SetParamNames("examId")
	c.SetParamValues(examId)

	s.mockExamService.EXPECT().
		FindExamRecords(int32(0), int32(10), examId, USER_ID).
		Return(&pb.FindExamRecordsResponse{
			Total:     1,
			PageCount: 1,
			ExamRecords: []*pb.ExamRecord{
				{
					Id:        "id01",
					ExamId:    examId,
					Score:     10,
					UserId:    USER_ID,
					CreatedAt: nil,
					UpdatedAt: nil,
				},
			},
		}, nil)

	// Test
	err := s.examHandler.FindExamRecords(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"total": 1, "pageCount": 1, "examRecords": [{
		"_id": "id01",
		"examId": "`+examId+`",
		"score": 10,
		"userId": "`+USER_ID+`",
		"createdAt": null,
		"updatedAt": null
	}]}`, rec.Body.String())
}

func (s *MyTestSuite) TestFindExamInfosWhenNotSignIn() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUserObjectId := primitive.NewObjectID()
	adminUserId := adminUserObjectId.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetAdminUser(context.TODO()).
		Return(&model.User{
			Id:       adminUserObjectId,
			Username: "TestAdmin",
			Password: "TestFindExamInfosWhenSignIn",
			Role:     "admin",
		}, nil)

	s.mockExamService.EXPECT().
		FindExamInfos(adminUserId, true).
		Return(&pb.FindExamInfosResponse{
			ExamInfos: []*pb.ExamInfo{
				{
					ExamId:        "exam01",
					Topic:         "t01",
					Description:   "d01",
					IsPublic:      true,
					QuestionCount: 10,
					RecordCount:   1,
				},
			},
		}, nil)

	// Test
	err := s.examHandler.FindExamInfosWhenNotSignIn(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"examInfos": [{
		"examId": "exam01",
		"topic": "t01",
		"description": "d01",
		"isPublic": true,
		"questionCount": 10,
		"recordCount": 1
	}]}`, rec.Body.String())
}

func (s *MyTestSuite) TestFindExamInfosWhenSignIn() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	adminUserObjectId := primitive.NewObjectID()
	adminUserId := adminUserObjectId.Hex()

	s.mockDatabaseRepository.EXPECT().
		GetAdminUser(context.TODO()).
		Return(&model.User{
			Id:       adminUserObjectId,
			Username: "TestAdmin",
			Password: "TestFindExamInfosWhenSignIn",
			Role:     "admin",
		}, nil)

	s.mockExamService.EXPECT().
		FindExamInfos(adminUserId, true).
		Return(&pb.FindExamInfosResponse{
			ExamInfos: []*pb.ExamInfo{
				{
					ExamId:        "exam01",
					Topic:         "t01",
					Description:   "d01",
					IsPublic:      true,
					QuestionCount: 10,
					RecordCount:   1,
				},
			},
		}, nil)
	s.mockExamService.EXPECT().
		FindExamInfos(USER_ID, true).
		Return(&pb.FindExamInfosResponse{
			ExamInfos: []*pb.ExamInfo{
				{
					ExamId:        "exam02",
					Topic:         "t02",
					Description:   "d02",
					IsPublic:      true,
					QuestionCount: 20,
					RecordCount:   2,
				},
			},
		}, nil)
	s.mockExamService.EXPECT().
		FindExamInfos(USER_ID, false).
		Return(&pb.FindExamInfosResponse{
			ExamInfos: []*pb.ExamInfo{
				{
					ExamId:        "exam03",
					Topic:         "t03",
					Description:   "d03",
					IsPublic:      false,
					QuestionCount: 30,
					RecordCount:   3,
				},
			},
		}, nil)

	// Test
	err := s.examHandler.FindExamInfosWhenSignIn(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"examInfos": [{
		"examId": "exam01",
		"topic": "t01",
		"description": "d01",
		"isPublic": true,
		"questionCount": 10,
		"recordCount": 1
	}, {
		"examId": "exam02",
		"topic": "t02",
		"description": "d02",
		"isPublic": true,
		"questionCount": 20,
		"recordCount": 2
	}, {
		"examId": "exam03",
		"topic": "t03",
		"description": "d03",
		"isPublic": false,
		"questionCount": 30,
		"recordCount": 3
	}]}`, rec.Body.String())
}
