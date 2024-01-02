package exam

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/examservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

const USER_ID = "user01"

type MyTestSuite struct {
	suite.Suite
	examHandler     examHandler
	mockExamService *examservice.MockExamService
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
		examServce: nil,
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
	s.examHandler.examServce = mockExamService
	s.mockExamService = mockExamService
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
