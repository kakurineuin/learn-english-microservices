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

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
)

// 使用測試的資料庫
const DATABASE = "learnEnglish_test"

type MyTestSuite struct {
	suite.Suite
	repo                          DatabaseRepository
	uri                           string
	ctx                           context.Context
	mongodbContainer              *mongodb.MongoDBContainer
	client                        *mongo.Client
	wordMeaningCollection         *mongo.Collection
	favoriteWordMeaningCollection *mongo.Collection
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
	s.wordMeaningCollection = client.Database(DATABASE).Collection("wordmeanings")
	s.favoriteWordMeaningCollection = client.Database(DATABASE).Collection("favoritewordmeanings")
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

func (s *MyTestSuite) TestCreateWordMeanings() {
	ctx := context.TODO()

	now := time.Now()
	size := 10
	wordMeanings := []model.WordMeaning{}

	for i := 0; i < size; i++ {
		wordMeanings = append(wordMeanings, model.WordMeaning{
			Word:         "test",
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
			QueryByWords:          []string{"test"},
			FavoriteWordMeaningId: primitive.NewObjectID(),
			CreatedAt:             now,
			UpdatedAt:             now,
		})
	}

	wordMeaningIds, err := s.repo.CreateWordMeanings(ctx, wordMeanings)
	s.Nil(err)
	s.Equal(size, len(wordMeaningIds))
}

func (s *MyTestSuite) TestFindWordMeaningsByWordAndUserId() {
	ctx := context.TODO()

	size := 10
	word := "TestFindWordMeaningsByWordAndUserId"
	userId := "user01"
	now := time.Now()
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

	wordMeanings, err := s.repo.FindWordMeaningsByWordAndUserId(ctx, word, userId)
	s.Nil(err)
	s.Equal(size, len(wordMeanings))
}

func (s *MyTestSuite) TestCreateFavoriteWordMeaning() {
	ctx := context.TODO()

	userId := "user01"
	wordMeaningId := "658dfc0d26c7337ddf4ab0cf"
	favoriteWordMeaningId, err := s.repo.CreateFavoriteWordMeaning(ctx, userId, wordMeaningId)
	s.Nil(err)
	s.NotEmpty(favoriteWordMeaningId)
}

func (s *MyTestSuite) TestGetFavoriteWordMeaningById() {
	ctx := context.TODO()

	userId := "user01"
	wordMeaningId := primitive.NewObjectID()
	result, err := s.favoriteWordMeaningCollection.InsertOne(ctx, model.FavoriteWordMeaning{
		UserId:        userId,
		WordMeaningId: wordMeaningId,
	})
	s.Nil(err)
	favoriteWordMeaningId := result.InsertedID.(primitive.ObjectID).Hex()

	favoriteWordMeaning, err := s.repo.GetFavoriteWordMeaningById(ctx, favoriteWordMeaningId)
	s.Nil(err)
	s.NotEmpty(favoriteWordMeaning)
}

func (s *MyTestSuite) TestFindFavoriteWordMeaningsByUserIdAndWord() {
	ctx := context.TODO()

	userId := "user01"
	word := "TestFindFavoriteWordMeaningsByUserIdAndWord"
	now := time.Now()
	size := 30
	wordMeaningDocuments := []interface{}{}

	for i := 0; i < size; i++ {
		wordMeaningDocuments = append(wordMeaningDocuments, model.WordMeaning{
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

	result, err := s.wordMeaningCollection.InsertMany(ctx, wordMeaningDocuments)
	s.Nil(err)
	s.Equal(size, len(result.InsertedIDs))

	favoriteWordMeaningDocuments := []interface{}{}

	for i := 0; i < size; i++ {
		favoriteWordMeaningDocuments = append(
			favoriteWordMeaningDocuments,
			model.FavoriteWordMeaning{
				UserId:        userId,
				WordMeaningId: result.InsertedIDs[i].(primitive.ObjectID),
			},
		)
	}

	result2, err := s.favoriteWordMeaningCollection.InsertMany(ctx, favoriteWordMeaningDocuments)
	s.Nil(err)
	s.Equal(size, len(result2.InsertedIDs))

	skip := int32(10)
	limit := int32(10)
	wordMeanings, err := s.repo.FindFavoriteWordMeaningsByUserIdAndWord(
		ctx,
		userId,
		word,
		skip,
		limit,
	)
	s.Nil(err)
	s.Equal(int(limit), len(wordMeanings))
}

func (s *MyTestSuite) TestCountFavoriteWordMeaningsByUserIdAndWord() {
	ctx := context.TODO()

	userId := "user01"
	word := "TestCountFavoriteWordMeaningsByUserIdAndWord"
	now := time.Now()
	size := 30
	wordMeaningDocuments := []interface{}{}

	for i := 0; i < size; i++ {
		wordMeaningDocuments = append(wordMeaningDocuments, model.WordMeaning{
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

	result, err := s.wordMeaningCollection.InsertMany(ctx, wordMeaningDocuments)
	s.Nil(err)
	s.Equal(size, len(result.InsertedIDs))

	favoriteWordMeaningDocuments := []interface{}{}

	for i := 0; i < size; i++ {
		favoriteWordMeaningDocuments = append(
			favoriteWordMeaningDocuments,
			model.FavoriteWordMeaning{
				UserId:        userId,
				WordMeaningId: result.InsertedIDs[i].(primitive.ObjectID),
			},
		)
	}

	result2, err := s.favoriteWordMeaningCollection.InsertMany(ctx, favoriteWordMeaningDocuments)
	s.Nil(err)
	s.Equal(size, len(result2.InsertedIDs))

	count, err := s.repo.CountFavoriteWordMeaningsByUserIdAndWord(
		ctx,
		userId,
		word,
	)
	s.Nil(err)
	s.Equal(int32(size), count)
}

func (s *MyTestSuite) TestDeleteFavoriteWordMeaningById() {
	ctx := context.TODO()

	userId := "user01"
	wordMeaningId := primitive.NewObjectID()
	result, err := s.favoriteWordMeaningCollection.InsertOne(ctx, model.FavoriteWordMeaning{
		UserId:        userId,
		WordMeaningId: wordMeaningId,
	})
	s.Nil(err)
	favoriteWordMeaningId := result.InsertedID.(primitive.ObjectID).Hex()

	deletedCount, err := s.repo.DeleteFavoriteWordMeaningById(ctx, favoriteWordMeaningId)
	s.Nil(err)
	s.Equal(int32(1), deletedCount)
}
