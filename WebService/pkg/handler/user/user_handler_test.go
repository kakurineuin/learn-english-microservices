package user

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/repository"
)

const TOKEN = "jwt_abc123"

type MyTestSuite struct {
	suite.Suite
	userHandler            userHandler
	mockDatabaseRepository *repository.MockDatabaseRepository
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}

// run once, before test suite methods
func (s *MyTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	// Mock JWT
	utilGetJWTToken = func(userId, username, role string) (string, error) {
		return TOKEN, nil
	}

	s.userHandler = userHandler{
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
	s.userHandler.databaseRepository = mockDatabaseRepository
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

func (s *MyTestSuite) TestCreateUser() {
	// Setup
	requestJSON := `{
  	"username": "someone",
  	"password": "TestCreateUser"
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userId := "user01"

	s.mockDatabaseRepository.EXPECT().
		GetUserByUsername(mock.Anything, "someone").
		Return(nil, nil)
	s.mockDatabaseRepository.EXPECT().
		CreateUser(mock.Anything, mock.Anything).
		Return(userId, nil)

	// Test
	err := s.userHandler.CreateUser(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"token": "`+TOKEN+`"}`, rec.Body.String())
}

func (s *MyTestSuite) TestLogin() {
	// Setup
	requestJSON := `{
  	"username": "someone",
  	"password": "TestLogin"
	}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockDatabaseRepository.EXPECT().
		GetUserByUsername(mock.Anything, "someone").
		Return(&model.User{
			Password: "$2b$10$jALDIZx8BaYAzljEiwXUZeOVniQIgvB20VJ.a4r94xjtwQB/eNIWa",
		}, nil)

	// Test
	err := s.userHandler.Login(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"token": "`+TOKEN+`"}`, rec.Body.String())
}
