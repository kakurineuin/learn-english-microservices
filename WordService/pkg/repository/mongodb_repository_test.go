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
	err = s.repo.ConnectDB(uri)
	s.Nil(err)

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
	type args struct {
		ctx          context.Context
		wordMeanings []model.WordMeaning
	}

	type setupDBResult struct{}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Create wordMeanings01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
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

				return &args{
					ctx:          ctx,
					wordMeanings: wordMeanings,
				}
			},
		},
		{
			name: "Create wordMeanings02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				now := time.Now()
				size := 3
				wordMeanings := []model.WordMeaning{}

				for i := 0; i < size; i++ {
					wordMeanings = append(wordMeanings, model.WordMeaning{
						Word:         "book",
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

				return &args{
					ctx:          ctx,
					wordMeanings: wordMeanings,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			wordMeaningIds, err := s.repo.CreateWordMeanings(
				args.ctx,
				args.wordMeanings,
			)
			s.Nil(err)
			s.Equal(len(args.wordMeanings), len(wordMeaningIds))
		})
	}
}

func (s *MyTestSuite) TestFindWordMeaningsByWordAndUserId() {
	type args struct {
		ctx    context.Context
		word   string
		userId string
	}

	type setupDBResult struct {
		word   string
		userId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name           string
		setupDB        func(s *MyTestSuite) *setupDBResult
		newArgs        func(dbResult setupDBResult) *args
		expectedLength int
	}{
		{
			name: "Find wordMeanings01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				size := 10
				word := "TestFindWordMeaningsByWordAndUserId01"
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

				return &setupDBResult{
					word:   word,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					word:   dbResult.word,
					userId: dbResult.userId,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find wordMeanings02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				size := 3
				word := "TestFindWordMeaningsByWordAndUserId02"
				userId := "user02"
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

				return &setupDBResult{
					word:   word,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					word:   dbResult.word,
					userId: dbResult.userId,
				}
			},
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			wordMeanings, err := s.repo.FindWordMeaningsByWordAndUserId(
				args.ctx, args.word, args.userId,
			)
			s.Nil(err)
			s.Len(wordMeanings, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestCreateFavoriteWordMeaning() {
	type args struct {
		ctx           context.Context
		userId        string
		wordMeaningId string
	}

	type setupDBResult struct{}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Create favoriteWordMeaning01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:           ctx,
					userId:        "user01",
					wordMeaningId: primitive.NewObjectID().Hex(),
				}
			},
		},
		{
			name: "Create favoriteWordMeaning02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:           ctx,
					userId:        "user02",
					wordMeaningId: primitive.NewObjectID().Hex(),
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			favoriteWordMeaningId, err := s.repo.CreateFavoriteWordMeaning(
				args.ctx,
				args.userId,
				args.wordMeaningId,
			)
			s.Nil(err)
			s.NotEmpty(favoriteWordMeaningId)
		})
	}
}

func (s *MyTestSuite) TestGetFavoriteWordMeaningById() {
	type args struct {
		ctx                   context.Context
		favoriteWordMeaningId string
	}

	type setupDBResult struct {
		favoriteWordMeaningId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Get favoriteWordMeaning01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user01"
				wordMeaningId := primitive.NewObjectID()
				result, err := s.favoriteWordMeaningCollection.InsertOne(
					ctx,
					model.FavoriteWordMeaning{
						UserId:        userId,
						WordMeaningId: wordMeaningId,
					},
				)
				s.Nil(err)

				return &setupDBResult{
					favoriteWordMeaningId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					favoriteWordMeaningId: dbResult.favoriteWordMeaningId,
				}
			},
		},
		{
			name: "Get favoriteWordMeaning02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user02"
				wordMeaningId := primitive.NewObjectID()
				result, err := s.favoriteWordMeaningCollection.InsertOne(
					ctx,
					model.FavoriteWordMeaning{
						UserId:        userId,
						WordMeaningId: wordMeaningId,
					},
				)
				s.Nil(err)

				return &setupDBResult{
					favoriteWordMeaningId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					favoriteWordMeaningId: dbResult.favoriteWordMeaningId,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			favoriteWordMeaning, err := s.repo.GetFavoriteWordMeaningById(
				args.ctx,
				args.favoriteWordMeaningId,
			)
			s.Nil(err)
			s.NotEmpty(favoriteWordMeaning)
		})
	}
}

func (s *MyTestSuite) TestFindFavoriteWordMeaningsByUserIdAndWord() {
	type args struct {
		ctx    context.Context
		userId string
		word   string
		skip   int32
		limit  int32
	}

	type setupDBResult struct {
		userId string
		word   string
	}

	ctx := context.TODO()

	testCases := []struct {
		name           string
		setupDB        func(s *MyTestSuite) *setupDBResult
		newArgs        func(dbResult setupDBResult) *args
		expectedLength int
	}{
		{
			name: "Find favoriteWordMeaning01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user01"
				word := "TestFindFavoriteWordMeaningsByUserIdAndWord01"
				now := time.Now()
				size := 10
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

				result2, err := s.favoriteWordMeaningCollection.InsertMany(
					ctx,
					favoriteWordMeaningDocuments,
				)
				s.Nil(err)
				s.Equal(size, len(result2.InsertedIDs))

				return &setupDBResult{
					userId: userId,
					word:   word,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					userId: dbResult.userId,
					word:   dbResult.word,
					skip:   0,
					limit:  10,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find favoriteWordMeaning02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user02"
				word := "TestFindFavoriteWordMeaningsByUserIdAndWord02"
				now := time.Now()
				size := 3
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

				result2, err := s.favoriteWordMeaningCollection.InsertMany(
					ctx,
					favoriteWordMeaningDocuments,
				)
				s.Nil(err)
				s.Equal(size, len(result2.InsertedIDs))

				return &setupDBResult{
					userId: userId,
					word:   word,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					userId: dbResult.userId,
					word:   dbResult.word,
					skip:   0,
					limit:  10,
				}
			},
			expectedLength: 3,
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			wordMeanings, err := s.repo.FindFavoriteWordMeaningsByUserIdAndWord(
				args.ctx,
				args.userId,
				args.word,
				args.skip,
				args.limit,
			)
			s.Nil(err)
			s.Len(wordMeanings, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestCountFavoriteWordMeaningsByUserIdAndWord() {
	type args struct {
		ctx    context.Context
		userId string
		word   string
	}

	type setupDBResult struct {
		userId string
		word   string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Count favoriteWordMeaning01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user01"
				word := "TestCountFavoriteWordMeaningsByUserIdAndWord01"
				now := time.Now()
				size := 10
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

				result2, err := s.favoriteWordMeaningCollection.InsertMany(
					ctx,
					favoriteWordMeaningDocuments,
				)
				s.Nil(err)
				s.Equal(size, len(result2.InsertedIDs))

				return &setupDBResult{
					userId: userId,
					word:   word,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					userId: dbResult.userId,
					word:   dbResult.word,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Count favoriteWordMeaning02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user02"
				word := "TestCountFavoriteWordMeaningsByUserIdAndWord02"
				now := time.Now()
				size := 3
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

				result2, err := s.favoriteWordMeaningCollection.InsertMany(
					ctx,
					favoriteWordMeaningDocuments,
				)
				s.Nil(err)
				s.Equal(size, len(result2.InsertedIDs))

				return &setupDBResult{
					userId: userId,
					word:   word,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					userId: dbResult.userId,
					word:   dbResult.word,
				}
			},
			expectedCount: 3,
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			count, err := s.repo.CountFavoriteWordMeaningsByUserIdAndWord(
				args.ctx,
				args.userId,
				args.word,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, count)
		})
	}
}

func (s *MyTestSuite) TestDeleteFavoriteWordMeaningById() {
	type args struct {
		ctx                   context.Context
		favoriteWordMeaningId string
	}

	type setupDBResult struct {
		favoriteWordMeaningId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Delete favoriteWordMeaning01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user01"
				wordMeaningId := primitive.NewObjectID()
				result, err := s.favoriteWordMeaningCollection.InsertOne(
					ctx,
					model.FavoriteWordMeaning{
						UserId:        userId,
						WordMeaningId: wordMeaningId,
					},
				)
				s.Nil(err)

				return &setupDBResult{
					favoriteWordMeaningId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:                   ctx,
					favoriteWordMeaningId: dbResult.favoriteWordMeaningId,
				}
			},
		},
		{
			name: "Delete favoriteWordMeaning02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user02"
				wordMeaningId := primitive.NewObjectID()
				result, err := s.favoriteWordMeaningCollection.InsertOne(
					ctx,
					model.FavoriteWordMeaning{
						UserId:        userId,
						WordMeaningId: wordMeaningId,
					},
				)
				s.Nil(err)

				return &setupDBResult{
					favoriteWordMeaningId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:                   ctx,
					favoriteWordMeaningId: dbResult.favoriteWordMeaningId,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			deletedCount, err := s.repo.DeleteFavoriteWordMeaningById(
				args.ctx,
				args.favoriteWordMeaningId,
			)
			s.Nil(err)
			s.EqualValues(1, deletedCount)
		})
	}
}
