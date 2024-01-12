package word

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/microservice/wordservice"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/util"
)

type MyIntegrationTestSuite struct {
	suite.Suite
	wordHandler           wordHandler
	compose               tc.ComposeStack
	client                *mongo.Client
	userCollection        *mongo.Collection
	wordMeaningCollection *mongo.Collection
	userId                string
	adminUserId           string
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
	err = compose.Up(ctx, tc.Wait(true))
	if err != nil {
		s.FailNow(err.Error())
	}

	databaseName := "Test_LearnEnglish"
	mongoDBURI := "mongodb://127.0.0.1:27017"

	// WordService
	wordService := wordservice.New(":8091")
	err = wordService.Connect()
	if err != nil {
		s.FailNow(err.Error())
	}

	// 用來建立測試資料的 client
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(mongoDBURI).SetTimeout(10*time.Second),
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.client = client
	s.userCollection = client.Database(databaseName).Collection("users")
	s.wordMeaningCollection = client.Database(databaseName).Collection("wordmeanings")

	// 新增測試資料
	now := time.Now()
	result, err := s.userCollection.InsertOne(context.TODO(), model.User{
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
	result, err = s.userCollection.InsertOne(context.TODO(), model.User{
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

	if err := s.client.Disconnect(context.TODO()); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// 程式結束時，結束微服務連線
	if err := s.wordHandler.wordService.Disconnect(); err != nil {
		log.Printf("wordService disconnect() error: %v", err)
	}

	// 終止 container
	if err := s.compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal); err != nil {
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
	ctx := context.TODO()
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
