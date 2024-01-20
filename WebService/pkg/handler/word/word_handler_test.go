package word

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
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/wordservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

type MyTestSuite struct {
	suite.Suite
	wordHandler     wordHandler
	mockWordService *wordservice.MockWordService
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
			UserId:           "user01",
			Username:         "test01",
			Role:             "user",
			RegisteredClaims: jwt.RegisteredClaims{},
		}
	}

	s.wordHandler = wordHandler{
		wordService: nil,
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
	mockWordService := wordservice.NewMockWordService(s.T())
	s.wordHandler.wordService = mockWordService
	s.mockWordService = mockWordService
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

func (s *MyTestSuite) TestFindWordMeanings() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("word")
	c.SetParamValues("test")

	s.mockWordService.EXPECT().
		FindWordByDictionary("test", "user01").
		Return(&pb.FindWordByDictionaryResponse{
			WordMeanings: []*pb.WordMeaning{},
		}, nil)

	// Test
	err := s.wordHandler.FindWordMeanings(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"wordMeanings": []}`, rec.Body.String())
}

func (s *MyTestSuite) TestCreateFavoriteWordMeaning() {
	// Setup
	requestJSON := `{
  	"wordMeaningId": "id01"
	}`
	e := echo.New()
	req := httptest.NewRequest(
		http.MethodPost,
		"/restricted/favorite",
		strings.NewReader(requestJSON),
	)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockWordService.EXPECT().
		CreateFavoriteWordMeaning("user01", "id01").
		Return(&pb.CreateFavoriteWordMeaningResponse{
			FavoriteWordMeaningId: "fid01",
		}, nil)

	// Test
	err := s.wordHandler.CreateFavoriteWordMeaning(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"favoriteWordMeaningId": "fid01"}`, rec.Body.String())
}

func (s *MyTestSuite) TestDeleteFavoriteWordMeaning() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("favoriteWordMeaningId")
	c.SetParamValues("f01")

	s.mockWordService.EXPECT().
		DeleteFavoriteWordMeaning("f01", "user01").
		Return(&pb.DeleteFavoriteWordMeaningResponse{}, nil)

	// Test
	err := s.wordHandler.DeleteFavoriteWordMeaning(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Empty(rec.Body.String())
}

func (s *MyTestSuite) TestFindFavoriteWordMeanings() {
	// Setup
	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	q.Set("word", "test")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockWordService.EXPECT().
		FindFavoriteWordMeanings(int32(0), int32(10), "user01", "test").
		Return(&pb.FindFavoriteWordMeaningsResponse{
			Total:                0,
			PageCount:            0,
			FavoriteWordMeanings: []*pb.WordMeaning{},
		}, nil)

	// Test
	err := s.wordHandler.FindFavoriteWordMeanings(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"total": 0, "pageCount": 0, "favoriteWordMeanings": []}`, rec.Body.String())
}

func (s *MyTestSuite) TestFindRandomFavoriteWordMeanings() {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.mockWordService.EXPECT().
		FindRandomFavoriteWordMeanings("user01", int32(10)).
		Return(&pb.FindRandomFavoriteWordMeaningsResponse{
			FavoriteWordMeanings: []*pb.WordMeaning{},
		}, nil)

	// Test
	err := s.wordHandler.FindRandomFavoriteWordMeanings(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.JSONEq(`{"favoriteWordMeanings": []}`, rec.Body.String())
}
