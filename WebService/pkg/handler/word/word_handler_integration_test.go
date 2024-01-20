package word

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

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/wordservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

type MyIntegrationTestSuite struct {
	suite.Suite
	wordHandler                   wordHandler
	compose                       tc.ComposeStack
	client                        *mongo.Client
	userCollection                *mongo.Collection
	wordMeaningCollection         *mongo.Collection
	favoriteWordMeaningCollection *mongo.Collection
	userId                        string
	adminUserId                   string
}

func TestMyIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(MyIntegrationTestSuite))
}

// run once, before test suite methods
func (s *MyIntegrationTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	// Setup WordService
	compose, err := tc.NewDockerCompose("./docker-compose.yml")
	if err != nil {
		s.FailNow(err.Error())
	}

	s.compose = compose

	ctx, cancel := context.WithCancel(context.Background())
	s.T().Cleanup(cancel)

	err = compose.
		WaitForService("word-service", wait.NewLogStrategy("Starting gRPC server at")).
		Up(ctx, tc.Wait(true))
	if err != nil {
		s.FailNow(err.Error())
	}

	databaseName := "Test_LearnEnglish"

	// 這裡改用 27018，避免同時跑 exam 和 word 兩隻整合測試時，兩個 mongo 搶佔同一個 port 的問題
	mongoDBURI := "mongodb://127.0.0.1:27018"

	// WordService
	wordService := wordservice.New(":8091")
	err = wordService.Connect()
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
	s.wordMeaningCollection = client.Database(databaseName).Collection("wordmeanings")
	s.favoriteWordMeaningCollection = client.Database(databaseName).
		Collection("favoritewordmeanings")

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

	s.wordHandler = wordHandler{
		wordService: wordService,
	}
}

// run once, after test suite methods
func (s *MyIntegrationTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
	ctx := context.Background()

	if err := s.client.Disconnect(ctx); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// 程式結束時，結束微服務連線
	if err := s.wordHandler.wordService.Disconnect(); err != nil {
		log.Printf("wordService disconnect() error: %v", err)
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

func (s *MyIntegrationTestSuite) TestFindWordMeanings() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	size := 10
	word := "test"
	documents := []interface{}{}

	for i := 0; i < size; i++ {
		documents = append(documents, model.WordMeaning{
			Word:         word,
			PartOfSpeech: "partOfSpeech",
			Gram:         "gram",
			Pronunciation: model.Pronunciation{
				Text:       "text",
				UkAudioUrl: "uk",
				UsAudioUrl: "us",
			},
			DefGram:    "defGram",
			Definition: fmt.Sprintf("this is a definition %d", i+1),
			Examples: []model.Example{
				{
					Pattern: "pattern",
					Examples: []model.Sentence{
						{
							AudioUrl: "audioUrl",
							Text:     "text",
						},
					},
				},
			},
			OrderByNo:             int32(i + 1),
			QueryByWords:          []string{word},
			FavoriteWordMeaningId: primitive.NewObjectID(),
			CreatedAt:             now,
			UpdatedAt:             now,
		})
	}

	_, err := s.wordMeaningCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("word")
	c.SetParamValues("test")

	// Test
	err = s.wordHandler.FindWordMeanings(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestCreateFavoriteWordMeaning() {
	// Setup
	requestJSON := `{
  	"wordMeaningId": "` + primitive.NewObjectID().Hex() + `"
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

	// Test
	err := s.wordHandler.CreateFavoriteWordMeaning(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestDeleteFavoriteWordMeaning() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	result, err := s.favoriteWordMeaningCollection.InsertOne(ctx, model.FavoriteWordMeaning{
		UserId:        s.userId,
		WordMeaningId: primitive.NewObjectID(),
		CreatedAt:     now,
		UpdatedAt:     now,
	})
	s.Nil(err)

	favoriteWordMeaningId := result.InsertedID.(primitive.ObjectID).Hex()

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("favoriteWordMeaningId")
	c.SetParamValues(favoriteWordMeaningId)

	// Test
	err = s.wordHandler.DeleteFavoriteWordMeaning(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Empty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindFavoriteWordMeanings() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	documents := []interface{}{}

	for i := 0; i < 10; i++ {
		documents = append(documents, model.FavoriteWordMeaning{
			UserId:        s.userId,
			WordMeaningId: primitive.NewObjectID(),
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}

	_, err := s.favoriteWordMeaningCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	q := make(url.Values)
	q.Set("pageIndex", "0")
	q.Set("pageSize", "10")
	q.Set("word", "test")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err = s.wordHandler.FindFavoriteWordMeanings(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}

func (s *MyIntegrationTestSuite) TestFindRandomFavoriteWordMeanings() {
	// Setup
	ctx := context.Background()
	now := time.Now()
	documents := []interface{}{}

	for i := 0; i < 10; i++ {
		documents = append(documents, model.FavoriteWordMeaning{
			UserId:        s.userId,
			WordMeaningId: primitive.NewObjectID(),
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}

	_, err := s.favoriteWordMeaningCollection.InsertMany(ctx, documents)
	s.Nil(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test
	err = s.wordHandler.FindRandomFavoriteWordMeanings(c)
	s.Nil(err)
	s.Equal(http.StatusOK, rec.Code)
	s.NotEmpty(rec.Body.String())
}
