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

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
)

// 使用測試的資料庫
const DATABASE = "learnEnglish_test"

type MyTestSuite struct {
	suite.Suite
	repo                  DatabaseRepository
	uri                   string
	ctx                   context.Context
	mongodbContainer      *mongodb.MongoDBContainer
	client                *mongo.Client
	examCollection        *mongo.Collection
	questionCollection    *mongo.Collection
	answerWrongCollection *mongo.Collection
	examRecordCollection  *mongo.Collection
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
		s.FailNow(err.Error())
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.repo = NewMongoDBRepository(DATABASE)
	err = s.repo.ConnectDB(uri)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.uri = uri
	s.ctx = ctx
	s.mongodbContainer = mongodbContainer

	// 用來建立測試資料的 client
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.client = client
	s.examCollection = client.Database(DATABASE).Collection("exams")
	s.questionCollection = client.Database(DATABASE).Collection("questions")
	s.answerWrongCollection = s.client.Database(DATABASE).Collection("answerwrongs")
	s.examRecordCollection = s.client.Database(DATABASE).Collection("examrecords")
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
		log.Printf("mongodbContainer.Terminate() error: %v", err)
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

func (s *MyTestSuite) TestCreateExam() {
	type args struct {
		ctx  context.Context
		exam model.Exam
	}

	type setupDBResult struct{}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Create exam01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx: ctx,
					exam: model.Exam{
						Topic:       "topic01",
						Description: "desc01",
						IsPublic:    true,
						Tags:        []string{"a01"},
						UserId:      "user01",
					},
				}
			},
		},
		{
			name: "Create exam02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx: ctx,
					exam: model.Exam{
						Topic:       "topic02",
						Description: "desc02",
						IsPublic:    false,
						Tags:        []string{"b01"},
						UserId:      "user02",
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			examId, err := s.repo.CreateExam(
				args.ctx,
				args.exam,
			)
			s.Nil(err)
			s.NotEmpty(examId)
		})
	}
}

func (s *MyTestSuite) TestUpdateExam() {
	type args struct {
		ctx  context.Context
		exam model.Exam
	}

	type setupDBResult struct {
		exam *model.Exam
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Update exam topic",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				exam := model.Exam{
					Topic:       "topic01",
					Description: "jsut for test",
					Tags:        []string{"tag01", "tag02"},
					IsPublic:    true,
					UserId:      "user03",
				}
				result, err := s.examCollection.InsertOne(ctx, exam)
				s.Nil(err)

				exam.Id = result.InsertedID.(primitive.ObjectID)

				return &setupDBResult{
					exam: &exam,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				exam := dbResult.exam
				exam.Topic = "TestCreateExam01-updated"
				return &args{
					ctx:  ctx,
					exam: *exam,
				}
			},
		},
		{
			name: "Update exam isPublic",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				exam := model.Exam{
					Topic:       "topic02",
					Description: "jsut for test",
					Tags:        []string{"tag01", "tag02"},
					IsPublic:    true,
					UserId:      "user04",
				}
				result, err := s.examCollection.InsertOne(ctx, exam)
				s.Nil(err)

				exam.Id = result.InsertedID.(primitive.ObjectID)

				return &setupDBResult{
					exam: &exam,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				exam := dbResult.exam
				exam.IsPublic = false
				return &args{
					ctx:  ctx,
					exam: *exam,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			err := s.repo.UpdateExam(
				args.ctx,
				args.exam,
			)
			s.Nil(err)
		})
	}
}

func (s *MyTestSuite) TestGetExamById() {
	type args struct {
		ctx    context.Context
		examId string
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Get exam01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				exam := model.Exam{
					Topic:       "topic01",
					Description: "jsut for test",
					Tags:        []string{"tag01", "tag02"},
					IsPublic:    true,
					UserId:      "user05",
				}
				result, err := s.examCollection.InsertOne(ctx, exam)
				s.Nil(err)

				return &setupDBResult{
					examId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
		},
		{
			name: "Get exam02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				exam := model.Exam{
					Topic:       "topic02",
					Description: "jsut for test",
					Tags:        []string{"tag01", "tag02"},
					IsPublic:    false,
					UserId:      "user06",
				}
				result, err := s.examCollection.InsertOne(ctx, exam)
				s.Nil(err)

				return &setupDBResult{
					examId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			exam, err := s.repo.GetExamById(
				args.ctx,
				args.examId,
			)
			s.Nil(err)
			s.NotNil(exam)
		})
	}
}

func (s *MyTestSuite) TestFindExamsByUserIdOrderByUpdateAtDesc() {
	type args struct {
		ctx    context.Context
		userId string
		skip   int32
		limit  int32
	}

	type setupDBResult struct {
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
			name: "Find exams by userId: user07",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user07"
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.Exam{
						Topic:       fmt.Sprintf("topic_%d", i),
						Description: "jsut for test",
						Tags:        []string{"tag01", "tag02"},
						IsPublic:    true,
						UserId:      userId,
					})
				}
				_, err := s.examCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					userId: dbResult.userId,
					skip:   0,
					limit:  10,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find exams by userId: user08",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user08"
				size := 13
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.Exam{
						Topic:       fmt.Sprintf("topic_%d", i),
						Description: "jsut for test",
						Tags:        []string{"tag01", "tag02"},
						IsPublic:    true,
						UserId:      userId,
					})
				}
				_, err := s.examCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					userId: dbResult.userId,
					skip:   10,
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
			exams, err := s.repo.FindExamsByUserIdOrderByUpdateAtDesc(
				args.ctx,
				args.userId,
				args.skip,
				args.limit,
			)
			s.Nil(err)
			s.Len(exams, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestFindExamsByUserIdAndIsPublicOrderByUpdateAtDesc() {
	type args struct {
		ctx      context.Context
		userId   string
		isPublic bool
	}

	type setupDBResult struct {
		userId   string
		isPublic bool
	}

	ctx := context.TODO()

	testCases := []struct {
		name           string
		setupDB        func(s *MyTestSuite) *setupDBResult
		newArgs        func(dbResult setupDBResult) *args
		expectedLength int
	}{
		{
			name: "Find exams by userId: user09, isPublic: true",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user09"
				isPublic := true
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.Exam{
						Topic:       fmt.Sprintf("topic_%d", i),
						Description: "jsut for test",
						Tags:        []string{"tag01", "tag02"},
						IsPublic:    isPublic,
						UserId:      userId,
					})
				}
				_, err := s.examCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					userId:   userId,
					isPublic: isPublic,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:      ctx,
					userId:   dbResult.userId,
					isPublic: dbResult.isPublic,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find exams by userId: user10, isPublic: false",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "user10"
				isPublic := false
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.Exam{
						Topic:       fmt.Sprintf("topic_%d", i),
						Description: "jsut for test",
						Tags:        []string{"tag01", "tag02"},
						IsPublic:    isPublic,
						UserId:      userId,
					})
				}
				_, err := s.examCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					userId:   userId,
					isPublic: isPublic,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:      ctx,
					userId:   dbResult.userId,
					isPublic: dbResult.isPublic,
				}
			},
			expectedLength: 10,
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			exams, err := s.repo.FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(
				args.ctx,
				args.userId,
				args.isPublic,
			)
			s.Nil(err)
			s.Len(exams, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestDeleteExamById() {
	type args struct {
		ctx    context.Context
		examId string
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Delete exam01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				exam := model.Exam{
					Topic:       "topic01",
					Description: "jsut for test",
					Tags:        []string{"tag01", "tag02"},
					IsPublic:    true,
					UserId:      "TestDeleteExamById",
				}
				result, err := s.examCollection.InsertOne(ctx, exam)
				s.Nil(err)

				return &setupDBResult{
					examId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
		},
		{
			name: "Delete exam02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				exam := model.Exam{
					Topic:       "topic02",
					Description: "desc02",
					Tags:        []string{},
					IsPublic:    false,
					UserId:      "TestDeleteExamById",
				}
				result, err := s.examCollection.InsertOne(ctx, exam)
				s.Nil(err)

				return &setupDBResult{
					examId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			deletedCount, err := s.repo.DeleteExamById(
				args.ctx,
				args.examId,
			)
			s.Nil(err)
			s.EqualValues(1, deletedCount)
		})
	}
}

func (s *MyTestSuite) TestCountExamsByUserId() {
	type args struct {
		ctx    context.Context
		userId string
	}

	type setupDBResult struct {
		userId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Count exams by userId: TestCountExamsByUserId01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "TestCountExamsByUserId01"
				documents := []interface{}{}
				size := 10

				for i := 0; i < size; i++ {
					documents = append(documents, model.Exam{
						Topic:       fmt.Sprintf("topic_%d", i),
						Description: "jsut for test",
						Tags:        []string{"tag01", "tag02"},
						IsPublic:    true,
						UserId:      userId,
					})
				}
				_, err := s.examCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					userId: dbResult.userId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Count exams by userId: TestCountExamsByUserId02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				userId := "TestCountExamsByUserId02"
				documents := []interface{}{}
				size := 13

				for i := 0; i < size; i++ {
					documents = append(documents, model.Exam{
						Topic:       fmt.Sprintf("topic_%d", i),
						Description: "desc02",
						Tags:        []string{},
						IsPublic:    false,
						UserId:      userId,
					})
				}
				_, err := s.examCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					userId: dbResult.userId,
				}
			},
			expectedCount: 13,
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			count, err := s.repo.CountExamsByUserId(
				args.ctx,
				args.userId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, count)
		})
	}
}

func (s *MyTestSuite) TestCreateQuestion() {
	type args struct {
		ctx      context.Context
		question model.Question
	}

	type setupDBResult struct{}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Create question01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx: ctx,
					question: model.Question{
						ExamId:  "exam01",
						Ask:     "q01",
						Answers: []string{"a01"},
						UserId:  "TestCreateQuestion01",
					},
				}
			},
		},
		{
			name: "Create question02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx: ctx,
					question: model.Question{
						ExamId:  "exam02",
						Ask:     "q02",
						Answers: []string{"a01", "a02"},
						UserId:  "TestCreateQuestion02",
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			questionId, err := s.repo.CreateQuestion(
				args.ctx,
				args.question,
			)
			s.Nil(err)
			s.NotEmpty(questionId)
		})
	}
}

func (s *MyTestSuite) TestUpdateQuestion() {
	type args struct {
		ctx      context.Context
		question model.Question
	}

	type setupDBResult struct {
		question *model.Question
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Update question ask",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				question := model.Question{
					ExamId:  "exam01",
					Ask:     "TestUpdateQuestion01",
					Answers: []string{"a01", "a02"},
					UserId:  "user01",
				}
				result, err := s.questionCollection.InsertOne(ctx, question)
				s.Nil(err)

				question.Id = result.InsertedID.(primitive.ObjectID)

				return &setupDBResult{
					question: &question,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				question := dbResult.question
				question.Ask = "TestUpdateQuestion01-updated"
				return &args{
					ctx:      ctx,
					question: *question,
				}
			},
		},
		{
			name: "Update question answers",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				question := model.Question{
					ExamId:  "exam01",
					Ask:     "TestUpdateQuestion02",
					Answers: []string{"a01"},
					UserId:  "user01",
				}
				result, err := s.questionCollection.InsertOne(ctx, question)
				s.Nil(err)

				question.Id = result.InsertedID.(primitive.ObjectID)

				return &setupDBResult{
					question: &question,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				question := dbResult.question
				question.Answers = []string{"b01", "b02"}
				return &args{
					ctx:      ctx,
					question: *question,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			err := s.repo.UpdateQuestion(
				args.ctx,
				args.question,
			)
			s.Nil(err)
		})
	}
}

func (s *MyTestSuite) TestGetQuestionById() {
	type args struct {
		ctx        context.Context
		questionId string
	}

	type setupDBResult struct {
		questionId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Get question01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				result, err := s.questionCollection.InsertOne(ctx, model.Question{
					ExamId:  "exam01",
					Ask:     "TestGetQuestion01",
					Answers: []string{"a01", "a02"},
					UserId:  "user01",
				})
				s.Nil(err)

				return &setupDBResult{
					questionId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					questionId: dbResult.questionId,
				}
			},
		},
		{
			name: "Get question02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				result, err := s.questionCollection.InsertOne(ctx, model.Question{
					ExamId:  "exam02",
					Ask:     "TestGetQuestion02",
					Answers: []string{"b01"},
					UserId:  "user02",
				})
				s.Nil(err)

				return &setupDBResult{
					questionId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					questionId: dbResult.questionId,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			question, err := s.repo.GetQuestionById(
				args.ctx,
				args.questionId,
			)
			s.Nil(err)
			s.NotNil(question)
		})
	}
}

func (s *MyTestSuite) TestFindQuestionsByQuestionIds() {
	type args struct {
		ctx         context.Context
		questionIds []string
	}

	type setupDBResult struct {
		questionIds []string
	}

	ctx := context.TODO()

	testCases := []struct {
		name           string
		setupDB        func(s *MyTestSuite) *setupDBResult
		newArgs        func(dbResult setupDBResult) *args
		expectedLength int
	}{
		{
			name: "Find questions01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindQuestionsByQuestionIds01"
				userId := "user01"
				documents := []interface{}{}
				size := 10

				for i := 0; i < size; i++ {
					documents = append(documents, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  userId,
					})
				}
				result, err := s.questionCollection.InsertMany(ctx, documents)
				s.Nil(err)

				questionIds := []string{}

				for _, id := range result.InsertedIDs {
					questionIds = append(questionIds, id.(primitive.ObjectID).Hex())
				}

				return &setupDBResult{
					questionIds: questionIds,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:         ctx,
					questionIds: dbResult.questionIds,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find questions02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindQuestionsByQuestionIds02"
				userId := "user01"
				documents := []interface{}{}
				size := 13

				for i := 0; i < size; i++ {
					documents = append(documents, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  userId,
					})
				}
				result, err := s.questionCollection.InsertMany(ctx, documents)
				s.Nil(err)

				questionIds := []string{}

				for _, id := range result.InsertedIDs {
					questionIds = append(questionIds, id.(primitive.ObjectID).Hex())
				}

				return &setupDBResult{
					questionIds: questionIds,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:         ctx,
					questionIds: dbResult.questionIds,
				}
			},
			expectedLength: 13,
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			questions, err := s.repo.FindQuestionsByQuestionIds(
				args.ctx,
				args.questionIds,
			)
			s.Nil(err)
			s.Len(questions, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestFindQuestionsByExamIdOrderByUpdateAtDesc() {
	type args struct {
		ctx    context.Context
		examId string
		skip   int32
		limit  int32
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name           string
		setupDB        func(s *MyTestSuite) *setupDBResult
		newArgs        func(dbResult setupDBResult) *args
		expectedLength int
	}{
		{
			name: "Find questions by examId: TestFindQuestionsByExamIdOrderByUpdateAtDesc01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindQuestionsByExamIdOrderByUpdateAtDesc01"
				userId := "user01"
				documents := []interface{}{}
				size := 10

				for i := 0; i < size; i++ {
					documents = append(documents, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  userId,
					})
				}
				_, err := s.questionCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					skip:   0,
					limit:  10,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find questions by examId: TestFindQuestionsByExamIdOrderByUpdateAtDesc02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindQuestionsByExamIdOrderByUpdateAtDesc02"
				userId := "user01"
				documents := []interface{}{}
				size := 3

				for i := 0; i < size; i++ {
					documents = append(documents, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  userId,
					})
				}
				_, err := s.questionCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
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
			questions, err := s.repo.FindQuestionsByExamIdOrderByUpdateAtDesc(
				args.ctx,
				args.examId,
				args.skip,
				args.limit,
			)
			s.Nil(err)
			s.Len(questions, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestDeleteQuestionById() {
	type args struct {
		ctx        context.Context
		questionId string
	}

	type setupDBResult struct {
		questionId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Delete question01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				result, err := s.questionCollection.InsertOne(ctx, model.Question{
					ExamId:  "exam01",
					Ask:     "TestDeleteQuestion01",
					Answers: []string{"a01", "a02"},
					UserId:  "user01",
				})
				s.Nil(err)

				return &setupDBResult{
					questionId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					questionId: dbResult.questionId,
				}
			},
		},
		{
			name: "Delete question02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				result, err := s.questionCollection.InsertOne(ctx, model.Question{
					ExamId:  "exam02",
					Ask:     "TestDeleteQuestion02",
					Answers: []string{"b01"},
					UserId:  "user02",
				})
				s.Nil(err)

				return &setupDBResult{
					questionId: result.InsertedID.(primitive.ObjectID).Hex(),
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					questionId: dbResult.questionId,
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			deletedCount, err := s.repo.DeleteQuestionById(
				args.ctx,
				args.questionId,
			)
			s.Nil(err)
			s.EqualValues(1, deletedCount)
		})
	}
}

func (s *MyTestSuite) TestDeleteQuestionsByExamId() {
	type args struct {
		ctx    context.Context
		examId string
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Delete questions by examId: TestDeleteQuestionsByExamId01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestDeleteQuestionsByExamId01"
				size := 10
				questions := []interface{}{}

				for i := 0; i < size; i++ {
					questions = append(questions, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  "user01",
					})
				}
				_, err := s.questionCollection.InsertMany(ctx, questions)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Delete questions by examId: TestDeleteQuestionsByExamId02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestDeleteQuestionsByExamId02"
				size := 3
				questions := []interface{}{}

				for i := 0; i < size; i++ {
					questions = append(questions, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  "user01",
					})
				}
				_, err := s.questionCollection.InsertMany(ctx, questions)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
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
			deletedCount, err := s.repo.DeleteQuestionsByExamId(
				args.ctx,
				args.examId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, deletedCount)
		})
	}
}

func (s *MyTestSuite) TestCountQuestionsByExamId() {
	type args struct {
		ctx    context.Context
		examId string
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Count questions by examId: TestCountQuestionsByExamId01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestCountQuestionsByExamId01"
				userId := "user01"
				size := 10
				questions := []interface{}{}

				for i := 0; i < size; i++ {
					questions = append(questions, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  userId,
					})
				}
				_, err := s.questionCollection.InsertMany(ctx, questions)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Count questions by examId: TestCountQuestionsByExamId02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestCountQuestionsByExamId02"
				userId := "user01"
				size := 3
				questions := []interface{}{}

				for i := 0; i < size; i++ {
					questions = append(questions, model.Question{
						ExamId:  examId,
						Ask:     fmt.Sprintf("Question_%d", i),
						Answers: []string{"a01", "a02"},
						UserId:  userId,
					})
				}
				_, err := s.questionCollection.InsertMany(ctx, questions)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
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
			count, err := s.repo.CountQuestionsByExamId(
				args.ctx,
				args.examId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, count)
		})
	}
}

func (s *MyTestSuite) TestDeleteAnswerWrongsByQuestionId() {
	type args struct {
		ctx        context.Context
		questionId string
	}

	type setupDBResult struct {
		questionId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Delete answerWrongs by questionId: TestDeleteAnswerWrongsByQuestionId01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				questionId := "TestDeleteAnswerWrongsByQuestionId01"
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.AnswerWrong{
						ExamId:     "exam_abc_01",
						QuestionId: questionId,
						Times:      10,
						UserId:     "user01",
					})
				}

				_, err := s.answerWrongCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					questionId: questionId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					questionId: dbResult.questionId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Delete answerWrongs by questionId: TestDeleteAnswerWrongsByQuestionId02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				questionId := "TestDeleteAnswerWrongsByQuestionId02"
				size := 3
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.AnswerWrong{
						ExamId:     "exam_abc_01",
						QuestionId: questionId,
						Times:      10,
						UserId:     "user01",
					})
				}

				_, err := s.answerWrongCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					questionId: questionId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					questionId: dbResult.questionId,
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
			deletedCount, err := s.repo.DeleteAnswerWrongsByQuestionId(
				args.ctx,
				args.questionId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, deletedCount)
		})
	}
}

func (s *MyTestSuite) TestDeleteAnswerWrongsByExamId() {
	type args struct {
		ctx    context.Context
		examId string
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Delete answerWrongs by examId: TestDeleteAnswerWrongsByExamId01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestDeleteAnswerWrongsByExamId01"
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.AnswerWrong{
						ExamId:     examId,
						QuestionId: "q01",
						Times:      10,
						UserId:     "user01",
					})
				}
				_, err := s.answerWrongCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					examId: dbResult.examId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Delete answerWrongs by examId: TestDeleteAnswerWrongsByExamId02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestDeleteAnswerWrongsByExamId02"
				size := 3
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.AnswerWrong{
						ExamId:     examId,
						QuestionId: "q01",
						Times:      10,
						UserId:     "user01",
					})
				}
				_, err := s.answerWrongCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					examId: dbResult.examId,
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
			deletedCount, err := s.repo.DeleteAnswerWrongsByExamId(
				args.ctx,
				args.examId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, deletedCount)
		})
	}
}

func (s *MyTestSuite) TestUpsertAnswerWrongByTimesPlusOne() {
	type args struct {
		ctx        context.Context
		examId     string
		questionId string
		userId     string
	}

	type setupDBResult struct{}

	type result struct {
		modifiedCount int32
		upsertedCount int32
	}

	ctx := context.TODO()

	testCases := []struct {
		name     string
		setupDB  func(s *MyTestSuite) *setupDBResult
		newArgs  func(dbResult setupDBResult) *args
		expected result
	}{
		{
			name: "Upsert(create) answerWrong by times plus one",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					examId:     "TestUpsertAnswerWrongByTimesPlusOne",
					questionId: "TestUpsertAnswerWrongByTimesPlusOne_q01",
					userId:     "TestUpsertAnswerWrongByTimesPlusOne_u01",
				}
			},
			expected: result{
				modifiedCount: 0,
				upsertedCount: 1,
			},
		},
		{
			name: "Upsert(update) answerWrong by times plus one",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:        ctx,
					examId:     "TestUpsertAnswerWrongByTimesPlusOne",
					questionId: "TestUpsertAnswerWrongByTimesPlusOne_q01",
					userId:     "TestUpsertAnswerWrongByTimesPlusOne_u01",
				}
			},
			expected: result{
				modifiedCount: 1,
				upsertedCount: 0,
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			modifiedCount, upsertedCount, err := s.repo.UpsertAnswerWrongByTimesPlusOne(
				args.ctx,
				args.examId,
				args.questionId,
				args.userId,
			)
			s.Nil(err)
			s.Equal(tc.expected.modifiedCount, modifiedCount)
			s.Equal(tc.expected.upsertedCount, upsertedCount)
		})
	}
}

func (s *MyTestSuite) TestFindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc() {
	type args struct {
		ctx    context.Context
		examId string
		userId string
		size   int32
	}

	type setupDBResult struct {
		examId string
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
			name: "Find answerWrongs by examId: FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc01"
				userId := "user01"
				documents := []interface{}{}
				size := 10

				for i := 0; i < size; i++ {
					documents = append(documents, model.AnswerWrong{
						ExamId:     examId,
						QuestionId: fmt.Sprintf("question%d", i+1),
						Times:      int32(i + 1),
						UserId:     userId,
					})
				}
				_, err := s.answerWrongCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					userId: dbResult.userId,
					size:   10,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find answerWrongs by examId: FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc02"
				userId := "user01"
				documents := []interface{}{}
				size := 3

				for i := 0; i < size; i++ {
					documents = append(documents, model.AnswerWrong{
						ExamId:     examId,
						QuestionId: fmt.Sprintf("question%d", i+1),
						Times:      int32(i + 1),
						UserId:     userId,
					})
				}
				_, err := s.answerWrongCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					userId: dbResult.userId,
					size:   10,
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
			answerWrongs, err := s.repo.FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc(
				args.ctx,
				args.examId,
				args.userId,
				args.size,
			)
			s.Nil(err)
			s.Len(answerWrongs, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestDeleteExamRecordsByExamId() {
	type args struct {
		ctx    context.Context
		examId string
	}

	type setupDBResult struct {
		examId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Delete examRecords by examId: TestDeleteExamRecordsByExamId01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestDeleteExamRecordsByExamId01"
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId: examId,
						Score:  6,
						UserId: "user01",
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Delete examRecords by examId: TestDeleteExamRecordsByExamId02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestDeleteExamRecordsByExamId02"
				size := 3
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId: examId,
						Score:  6,
						UserId: "user01",
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
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
			deletedCount, err := s.repo.DeleteExamRecordsByExamId(
				args.ctx,
				args.examId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, deletedCount)
		})
	}
}

func (s *MyTestSuite) TestCreateExamRecord() {
	type args struct {
		ctx        context.Context
		examRecord model.ExamRecord
	}

	type setupDBResult struct{}

	ctx := context.TODO()

	testCases := []struct {
		name    string
		setupDB func(s *MyTestSuite) *setupDBResult
		newArgs func(dbResult setupDBResult) *args
	}{
		{
			name: "Create examRecord01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx: ctx,
					examRecord: model.ExamRecord{
						ExamId: "TestCreateExamRecord01",
						Score:  10,
						UserId: "user01",
					},
				}
			},
		},
		{
			name: "Create examRecord02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				return &setupDBResult{}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx: ctx,
					examRecord: model.ExamRecord{
						ExamId: "TestCreateExamRecord02",
						Score:  6,
						UserId: "user02",
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.newArgs(*tc.setupDB(s))

			// Test
			examRecordId, err := s.repo.CreateExamRecord(
				args.ctx,
				args.examRecord,
			)
			s.Nil(err)
			s.NotEmpty(examRecordId)
		})
	}
}

func (s *MyTestSuite) TestFindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc() {
	type args struct {
		ctx    context.Context
		examId string
		userId string
		skip   int32
		limit  int32
	}

	type setupDBResult struct {
		examId string
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
			name: "Find examRecords 01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc01"
				userId := "user01"
				documents := []interface{}{}
				size := 10

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId: examId,
						Score:  10,
						UserId: userId,
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					userId: dbResult.userId,
					skip:   0,
					limit:  10,
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find examRecords 02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc02"
				userId := "user02"
				documents := []interface{}{}
				size := 3

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId: examId,
						Score:  10,
						UserId: userId,
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					userId: dbResult.userId,
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
			examRecords, err := s.repo.FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
				args.ctx,
				args.examId,
				args.userId,
				args.skip,
				args.limit,
			)
			s.Nil(err)
			s.Len(examRecords, tc.expectedLength)
		})
	}
}

func (s *MyTestSuite) TestCountExamRecordsByExamIdAndUserId() {
	type args struct {
		ctx    context.Context
		examId string
		userId string
	}

	type setupDBResult struct {
		examId string
		userId string
	}

	ctx := context.TODO()

	testCases := []struct {
		name          string
		setupDB       func(s *MyTestSuite) *setupDBResult
		newArgs       func(dbResult setupDBResult) *args
		expectedCount int32
	}{
		{
			name: "Count examRecords01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestCountExamRecordsByExamIdAndUserId01"
				userId := "user01"
				size := 10
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId: examId,
						Score:  10,
						UserId: userId,
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					userId: dbResult.userId,
				}
			},
			expectedCount: 10,
		},
		{
			name: "Count examRecords02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestCountExamRecordsByExamIdAndUserId02"
				userId := "user01"
				size := 3
				documents := []interface{}{}

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId: examId,
						Score:  10,
						UserId: userId,
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					ctx:    ctx,
					examId: dbResult.examId,
					userId: dbResult.userId,
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
			count, err := s.repo.CountExamRecordsByExamIdAndUserId(
				args.ctx,
				args.examId,
				args.userId,
			)
			s.Nil(err)
			s.Equal(tc.expectedCount, count)
		})
	}
}

func (s *MyTestSuite) TestFindExamRecordsByExamIdAndUserIdAndCreatedAt() {
	type args struct {
		ctx       context.Context
		examId    string
		userId    string
		createdAt time.Time
	}

	type setupDBResult struct {
		examId string
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
			name: "Find examRecords01",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindExamRecordsByExamIdAndUserIdAndCreatedAt01"
				userId := "user01"
				documents := []interface{}{}
				size := 10
				now := time.Now()

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId:    examId,
						Score:     10,
						UserId:    userId,
						CreatedAt: now,
						UpdatedAt: now,
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					examId:    dbResult.examId,
					userId:    dbResult.userId,
					createdAt: time.Now().AddDate(0, 0, -29),
				}
			},
			expectedLength: 10,
		},
		{
			name: "Find examRecords02",
			setupDB: func(s *MyTestSuite) *setupDBResult {
				examId := "TestFindExamRecordsByExamIdAndUserIdAndCreatedAt02"
				userId := "user02"
				documents := []interface{}{}
				size := 3
				now := time.Now()

				for i := 0; i < size; i++ {
					documents = append(documents, model.ExamRecord{
						ExamId:    examId,
						Score:     10,
						UserId:    userId,
						CreatedAt: now,
						UpdatedAt: now,
					})
				}
				_, err := s.examRecordCollection.InsertMany(ctx, documents)
				s.Nil(err)

				return &setupDBResult{
					examId: examId,
					userId: userId,
				}
			},
			newArgs: func(dbResult setupDBResult) *args {
				return &args{
					examId:    dbResult.examId,
					userId:    dbResult.userId,
					createdAt: time.Now().AddDate(0, 0, -29),
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
			examRecords, err := s.repo.FindExamRecordsByExamIdAndUserIdAndCreatedAt(
				args.ctx,
				args.examId,
				args.userId,
				args.createdAt,
			)
			s.Nil(err)
			s.Len(examRecords, tc.expectedLength)
		})
	}
}
