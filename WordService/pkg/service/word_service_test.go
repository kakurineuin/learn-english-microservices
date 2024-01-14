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
	s.wordService = wordService{
		logger:             logger,
		errorLogger:        level.Error(logger),
		databaseRepository: nil,
		spider:             nil,
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
	s.wordService.databaseRepository = mockDatabaseRepository
	s.mockDatabaseRepository = mockDatabaseRepository

	mockSpider := crawler.NewMockSpider(s.T())
	s.wordService.spider = mockSpider
	s.mockSpider = mockSpider
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
	type args struct {
		word   string
		userId string
	}

	type result struct {
		wordMeanings []model.WordMeaning
		err          error
	}

	word01 := "test"
	mockWordMeanings01 := []model.WordMeaning{
		{
			Word: word01,
		},
		{
			Word: word01,
		},
		{
			Word: word01,
		},
	}

	word02 := "book"
	mockWordMeanings02 := []model.WordMeaning{
		{
			Word: word02,
		},
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find wordMeanings01",
			args: &args{
				word:   word01,
				userId: "user01",
			},
			expected: &result{
				wordMeanings: mockWordMeanings01,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindWordMeaningsByWordAndUserId(
						mock.Anything, args.word, args.userId).
					Return(mockWordMeanings01, nil)
			},
		},
		{
			name: "Find wordMeanings02",
			args: &args{
				word:   word02,
				userId: "user01",
			},
			expected: &result{
				wordMeanings: mockWordMeanings02,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindWordMeaningsByWordAndUserId(
						mock.Anything, args.word, args.userId).
					Return(mockWordMeanings02, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			wordMeanings, err := s.wordService.FindWordByDictionary(
				args.word, args.userId,
			)
			expected := tc.expected
			s.Equal(expected.wordMeanings, wordMeanings)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindWordByDictionary_WhenDataFromCrawler() {
	type args struct {
		word   string
		userId string
	}

	type result struct {
		wordMeanings []model.WordMeaning
		err          error
	}

	word01 := "test"
	mockWordMeanings01 := []model.WordMeaning{
		{
			Word: word01,
		},
		{
			Word: word01,
		},
		{
			Word: word01,
		},
	}

	word02 := "book"
	mockWordMeanings02 := []model.WordMeaning{
		{
			Word: word02,
		},
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find wordMeanings01",
			args: &args{
				word:   word01,
				userId: "user01",
			},
			expected: &result{
				wordMeanings: mockWordMeanings01,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindWordMeaningsByWordAndUserId(
						mock.Anything, args.word, args.userId).
					Return(nil, nil).
					Once()
				s.mockSpider.EXPECT().FindWordMeaningsFromDictionary(args.word).
					Return(mockWordMeanings01, nil)
				s.mockDatabaseRepository.EXPECT().
					CreateWordMeanings(mock.Anything, mockWordMeanings01).
					Return([]string{"id1", "id2", "id3"}, nil)
				s.mockDatabaseRepository.EXPECT().
					FindWordMeaningsByWordAndUserId(
						mock.Anything, args.word, args.userId).
					Return(mockWordMeanings01, nil)
			},
		},
		{
			name: "Find wordMeanings02",
			args: &args{
				word:   word02,
				userId: "user01",
			},
			expected: &result{
				wordMeanings: mockWordMeanings02,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindWordMeaningsByWordAndUserId(
						mock.Anything, args.word, args.userId).
					Return(nil, nil).
					Once()
				s.mockSpider.EXPECT().FindWordMeaningsFromDictionary(args.word).
					Return(mockWordMeanings02, nil)
				s.mockDatabaseRepository.EXPECT().
					CreateWordMeanings(mock.Anything, mockWordMeanings02).
					Return([]string{"id1", "id2", "id3"}, nil)
				s.mockDatabaseRepository.EXPECT().
					FindWordMeaningsByWordAndUserId(
						mock.Anything, args.word, args.userId).
					Return(mockWordMeanings02, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			wordMeanings, err := s.wordService.FindWordByDictionary(
				args.word, args.userId,
			)
			expected := tc.expected
			s.Equal(expected.wordMeanings, wordMeanings)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestCreateFavoriteWordMeaning() {
	type args struct {
		userId        string
		wordMeaningId string
	}

	type result struct {
		favoriteWordMeaningId string
		err                   error
	}

	mockFavoriteWordMeaningId01 := "f01"
	mockFavoriteWordMeaningId02 := "f02"

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Create favoriteWordMeaning01",
			args: &args{
				userId:        "user01",
				wordMeaningId: "w01",
			},
			expected: &result{
				favoriteWordMeaningId: mockFavoriteWordMeaningId01,
				err:                   nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					CreateFavoriteWordMeaning(
						mock.Anything, args.userId, args.wordMeaningId).
					Return(mockFavoriteWordMeaningId01, nil)
			},
		},
		{
			name: "Create favoriteWordMeaning02",
			args: &args{
				userId:        "user02",
				wordMeaningId: "w02",
			},
			expected: &result{
				favoriteWordMeaningId: mockFavoriteWordMeaningId02,
				err:                   nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					CreateFavoriteWordMeaning(
						mock.Anything, args.userId, args.wordMeaningId).
					Return(mockFavoriteWordMeaningId02, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			favoriteWordMeaningId, err := s.wordService.CreateFavoriteWordMeaning(
				args.userId, args.wordMeaningId,
			)
			expected := tc.expected
			s.Equal(expected.favoriteWordMeaningId, favoriteWordMeaningId)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestDeleteFavoriteWordMeaning() {
	type args struct {
		favoriteWordMeaningId string
		userId                string
	}

	type result struct {
		err error
	}

	favoriteWordMeaningId01 := primitive.NewObjectID().Hex()
	favoriteWordMeaningId02 := primitive.NewObjectID().Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Delete favoriteWordMeaning01",
			args: &args{
				favoriteWordMeaningId: favoriteWordMeaningId01,
				userId:                "user01",
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetFavoriteWordMeaningById(
						mock.Anything, args.favoriteWordMeaningId).
					Return(&model.FavoriteWordMeaning{
						UserId:        args.userId,
						WordMeaningId: primitive.NewObjectID(),
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					DeleteFavoriteWordMeaningById(
						mock.Anything, args.favoriteWordMeaningId).
					Return(int32(1), nil)
			},
		},
		{
			name: "Delete favoriteWordMeaning02",
			args: &args{
				favoriteWordMeaningId: favoriteWordMeaningId02,
				userId:                "user02",
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetFavoriteWordMeaningById(
						mock.Anything, args.favoriteWordMeaningId).
					Return(&model.FavoriteWordMeaning{
						UserId:        args.userId,
						WordMeaningId: primitive.NewObjectID(),
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					DeleteFavoriteWordMeaningById(
						mock.Anything, args.favoriteWordMeaningId).
					Return(int32(1), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			err := s.wordService.DeleteFavoriteWordMeaning(
				args.favoriteWordMeaningId,
				args.userId,
			)
			expected := tc.expected
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindFavoriteWordMeanings() {
	type args struct {
		pageIndex int32
		pageSize  int32
		userId    string
		word      string
	}

	type result struct {
		total        int32
		pageCount    int32
		wordMeanings []model.WordMeaning
		err          error
	}

	word01 := "test"
	mockWordMeanings01 := []model.WordMeaning{
		{
			Word: word01,
		},
		{
			Word: word01,
		},
		{
			Word: word01,
		},
	}

	word02 := "book"
	mockWordMeanings02 := []model.WordMeaning{
		{
			Word: word02,
		},
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find favoriteWordMeanings01",
			args: &args{
				pageIndex: 0,
				pageSize:  10,
				userId:    "user01",
				word:      word01,
			},
			expected: &result{
				total:        3,
				pageCount:    1,
				wordMeanings: mockWordMeanings01,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				skip := args.pageSize * args.pageIndex
				limit := args.pageSize

				s.mockDatabaseRepository.EXPECT().
					FindFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything, args.userId, args.word, skip, limit).
					Return(mockWordMeanings01, nil)
				s.mockDatabaseRepository.EXPECT().
					CountFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything, args.userId, args.word).
					Return(int32(3), nil)
			},
		},
		{
			name: "Find favoriteWordMeanings02",
			args: &args{
				pageIndex: 1,
				pageSize:  10,
				userId:    "user02",
				word:      word02,
			},
			expected: &result{
				total:        11,
				pageCount:    2,
				wordMeanings: mockWordMeanings02,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				skip := args.pageSize * args.pageIndex
				limit := args.pageSize

				s.mockDatabaseRepository.EXPECT().
					FindFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything, args.userId, args.word, skip, limit).
					Return(mockWordMeanings02, nil)
				s.mockDatabaseRepository.EXPECT().
					CountFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything, args.userId, args.word).
					Return(int32(11), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			total, pageCount, wordMeanings, err := s.wordService.FindFavoriteWordMeanings(
				args.pageIndex,
				args.pageSize,
				args.userId,
				args.word,
			)
			expected := tc.expected
			s.Equal(expected.total, total)
			s.Equal(expected.pageCount, pageCount)
			s.Equal(expected.wordMeanings, wordMeanings)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindRandomFavoriteWordMeanings() {
	type args struct {
		userId string
		size   int32
	}

	type result struct {
		wordMeanings []model.WordMeaning
		err          error
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find random favoriteWordMeanings01",
			args: &args{
				userId: "user01",
				size:   10,
			},
			expected: &result{
				wordMeanings: []model.WordMeaning{
					{}, {}, {},
				},
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything,
						args.userId,
						"",
						mock.AnythingOfType("int32"),
						int32(1),
					).
					Return([]model.WordMeaning{
						{},
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CountFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything, args.userId, "").
					Return(int32(3), nil)
			},
		},
		{
			name: "Find random favoriteWordMeanings02",
			args: &args{
				userId: "user02",
				size:   10,
			},
			expected: &result{
				wordMeanings: []model.WordMeaning{
					{},
				},
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything,
						args.userId,
						"",
						mock.AnythingOfType("int32"),
						int32(1),
					).
					Return([]model.WordMeaning{
						{},
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CountFavoriteWordMeaningsByUserIdAndWord(
						mock.Anything, args.userId, "").
					Return(int32(1), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			wordMeanings, err := s.wordService.FindRandomFavoriteWordMeanings(
				args.userId,
				args.size,
			)
			expected := tc.expected
			s.Equal(expected.wordMeanings, wordMeanings)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestCircleCI() {
	// test 7
	s.Fail("========= WordService For test CircleCI !!!")
}
