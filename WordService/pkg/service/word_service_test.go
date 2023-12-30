package service

import (
	"log"
	"os"
	"testing"

	gokitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/crawler"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/word-service/pkg/repository"
)

type MyTestSuite struct {
	suite.Suite
	wordService            wordService
	mockDatabaseRepository *repository.MockDatabaseRepository
	mockSpider             *crawler.MockSpider
}

func TestMyTestSuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}

// run once, before test suite methods
func (s *MyTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	logger := gokitlog.NewJSONLogger(os.Stdout)
	logger = gokitlog.With(
		logger,
		"ts",
		gokitlog.DefaultTimestampUTC,
		"caller",
		gokitlog.DefaultCaller,
	)
	mockDatabaseRepository := repository.NewMockDatabaseRepository(s.T())
	mockSpider := crawler.NewMockSpider(s.T())
	s.wordService = wordService{
		logger:             logger,
		errorLogger:        level.Error(logger),
		databaseRepository: mockDatabaseRepository,
		spider:             mockSpider,
	}
	s.mockDatabaseRepository = mockDatabaseRepository
	s.mockSpider = mockSpider
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
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

func (s *MyTestSuite) TestFindWordByDictionary_WhenDataFromDB() {
	word := "test"
	userId := "user01"
	mockWordMeanings := []model.WordMeaning{
		{
			Word: word,
		},
		{
			Word: word,
		},
		{
			Word: word,
		},
	}
	size := len(mockWordMeanings)

	s.mockDatabaseRepository.EXPECT().
		FindWordMeaningsByWordAndUserId(mock.Anything, word, userId).
		Return(mockWordMeanings, nil)

	wordMeanings, err := s.wordService.FindWordByDictionary(word, userId)
	s.Nil(err)
	s.Equal(size, len(wordMeanings))
}

func (s *MyTestSuite) TestFindWordByDictionary_WhenDataFromCrawler() {
	word := "test"
	userId := "user01"
	mockWordMeanings := []model.WordMeaning{
		{
			Word: word,
		},
		{
			Word: word,
		},
		{
			Word: word,
		},
	}
	size := len(mockWordMeanings)

	s.mockDatabaseRepository.EXPECT().
		FindWordMeaningsByWordAndUserId(mock.Anything, word, userId).
		Return(nil, nil).Once()
	s.mockSpider.EXPECT().FindWordMeaningsFromDictionary(word).Return(mockWordMeanings, nil)
	s.mockDatabaseRepository.EXPECT().
		CreateWordMeanings(mock.Anything, mockWordMeanings).
		Return([]string{"id1", "id2", "id3"}, nil)
	s.mockDatabaseRepository.EXPECT().
		FindWordMeaningsByWordAndUserId(mock.Anything, word, userId).
		Return(mockWordMeanings, nil)

	wordMeanings, err := s.wordService.FindWordByDictionary(word, userId)
	s.Nil(err)
	s.Equal(size, len(wordMeanings))
}

func (s *MyTestSuite) TestCreateFavoriteWordMeaning() {
	userId := "user01"
	wordMeaningId := "aaa01"
	mockFavoriteWordMeaningId := "bbb01"

	s.mockDatabaseRepository.EXPECT().
		CreateFavoriteWordMeaning(mock.Anything, userId, wordMeaningId).
		Return(mockFavoriteWordMeaningId, nil)

	favoriteWordMeaningId, err := s.wordService.CreateFavoriteWordMeaning(userId, wordMeaningId)
	s.Nil(err)
	s.Equal(mockFavoriteWordMeaningId, favoriteWordMeaningId)
}

func (s *MyTestSuite) TestDeleteFavoriteWordMeaning() {
	userId := "user01"
	favoriteWordMeaningId := primitive.NewObjectID().Hex()

	s.mockDatabaseRepository.EXPECT().
		GetFavoriteWordMeaningById(mock.Anything, favoriteWordMeaningId).
		Return(&model.FavoriteWordMeaning{
			UserId:        userId,
			WordMeaningId: primitive.NewObjectID(),
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		DeleteFavoriteWordMeaningById(mock.Anything, favoriteWordMeaningId).
		Return(int64(1), nil)

	err := s.wordService.DeleteFavoriteWordMeaning(favoriteWordMeaningId, userId)
	s.Nil(err)
}

func (s *MyTestSuite) TestFindFavoriteWordMeanings() {
	userId := "user01"
	word := "TestFindFavoriteWordMeanings"
	pageIndex := int64(1)
	pageSize := int64(10)
	skip := pageSize * pageIndex
	limit := pageSize
	mockTotal := int64(13)

	s.mockDatabaseRepository.EXPECT().
		FindFavoriteWordMeaningsByUserIdAndWord(mock.Anything, userId, word, skip, limit).
		Return([]model.WordMeaning{
			{},
			{},
			{},
		}, nil)
	s.mockDatabaseRepository.EXPECT().
		CountFavoriteWordMeaningsByUserIdAndWord(mock.Anything, userId, word).
		Return(mockTotal, nil)

	total, pageCount, wordMeanings, err := s.wordService.FindFavoriteWordMeanings(
		pageIndex,
		pageSize,
		userId,
		word,
	)
	s.Nil(err)
	s.Equal(total, mockTotal)
	s.Equal(int64(2), pageCount)
	s.Equal(3, len(wordMeanings))
}
