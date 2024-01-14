package service

import (
	"log"
	"os"
	"testing"
	"time"

	gokitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/model"
	"github.com/kakurineuin/learn-english-microservices/exam-service/pkg/repository"
)

type MyTestSuite struct {
	suite.Suite
	examService            examService
	mockDatabaseRepository *repository.MockDatabaseRepository
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
	s.examService = examService{
		logger:             logger,
		errorLogger:        level.Error(logger),
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
	s.examService.databaseRepository = mockDatabaseRepository
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

func (s *MyTestSuite) TestCreateExam() {
	type args struct {
		topic       string
		description string
		isPublic    bool
		userId      string
	}

	type result struct {
		examId string
		err    error
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Create exam01",
			args: &args{
				topic:       "topic01",
				description: "desc01",
				isPublic:    true,
				userId:      "user01",
			},
			expected: &result{
				examId: "exam01",
				err:    nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					CreateExam(mock.Anything, model.Exam{
						Topic:       args.topic,
						Description: args.description,
						Tags:        []string{},
						IsPublic:    args.isPublic,
						UserId:      args.userId,
					}).
					Return("exam01", nil)
			},
		},
		{
			name: "Create exam02",
			args: &args{
				topic:       "topic02",
				description: "desc02",
				isPublic:    false,
				userId:      "user02",
			},
			expected: &result{
				examId: "exam02",
				err:    nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					CreateExam(mock.Anything, model.Exam{
						Topic:       args.topic,
						Description: args.description,
						Tags:        []string{},
						IsPublic:    args.isPublic,
						UserId:      args.userId,
					}).
					Return("exam02", nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			examId, err := s.examService.CreateExam(
				args.topic,
				args.description,
				args.isPublic,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.examId, examId)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestUpdateExam() {
	type args struct {
		examId      string
		topic       string
		description string
		isPublic    bool
		userId      string
	}

	type result struct {
		examId string
		err    error
	}

	userId := "user01"
	id := primitive.NewObjectID()
	examId := id.Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Update exam topic and description",
			args: &args{
				examId:      examId,
				topic:       "topic01-u01",
				description: "desc01-u01",
				isPublic:    true,
				userId:      userId,
			},
			expected: &result{
				examId: examId,
				err:    nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id,
						Topic:       "topic01",
						Description: "desc01",
						Tags:        []string{},
						IsPublic:    true,
						UserId:      userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					UpdateExam(mock.Anything, model.Exam{
						Id:          id,
						Topic:       args.topic,
						Description: args.description,
						Tags:        []string{},
						IsPublic:    args.isPublic,
						UserId:      args.userId,
					}).
					Return(nil)
			},
		},
		{
			name: "Update exam isPublic",
			args: &args{
				examId:      examId,
				topic:       "topic01",
				description: "desc01",
				isPublic:    false,
				userId:      userId,
			},
			expected: &result{
				examId: examId,
				err:    nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id,
						Topic:       "topic01",
						Description: "desc01",
						Tags:        []string{},
						IsPublic:    true,
						UserId:      userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					UpdateExam(mock.Anything, model.Exam{
						Id:          id,
						Topic:       args.topic,
						Description: args.description,
						Tags:        []string{},
						IsPublic:    args.isPublic,
						UserId:      args.userId,
					}).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			examId, err := s.examService.UpdateExam(
				args.examId,
				args.topic,
				args.description,
				args.isPublic,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.examId, examId)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindExams() {
	type args struct {
		pageIndex int32
		pageSize  int32
		userId    string
	}

	type result struct {
		total     int32
		pageCount int32
		exams     []model.Exam
		err       error
	}

	userId := "user01"
	mockExams := []model.Exam{
		{},
		{},
		{},
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find exams by page1",
			args: &args{
				pageIndex: 0,
				pageSize:  10,
				userId:    userId,
			},
			expected: &result{
				total:     3,
				pageCount: 1,
				exams:     mockExams,
				err:       nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindExamsByUserIdOrderByUpdateAtDesc(
						mock.Anything,
						args.userId,
						args.pageIndex*args.pageSize,
						args.pageSize,
					).
					Return(mockExams, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamsByUserId(mock.Anything, args.userId).
					Return(int32(3), nil)
			},
		},
		{
			name: "Find exams by page2",
			args: &args{
				pageIndex: 1,
				pageSize:  10,
				userId:    userId,
			},
			expected: &result{
				total:     13,
				pageCount: 2,
				exams:     mockExams,
				err:       nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindExamsByUserIdOrderByUpdateAtDesc(
						mock.Anything,
						args.userId,
						args.pageIndex*args.pageSize,
						args.pageSize,
					).
					Return(mockExams, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamsByUserId(mock.Anything, args.userId).
					Return(int32(13), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			total, pageCount, exams, err := s.examService.FindExams(
				args.pageIndex,
				args.pageSize,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.total, total)
			s.Equal(expected.pageCount, pageCount)
			s.Equal(expected.exams, exams)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestDeleteExam() {
	type args struct {
		examId string
		userId string
	}

	type result struct {
		err error
	}

	userId := "user01"
	id01 := primitive.NewObjectID()
	examId01 := id01.Hex()

	id02 := primitive.NewObjectID()
	examId02 := id02.Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Delete exam01",
			args: &args{
				examId: examId01,
				userId: userId,
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id01,
						Topic:       "topic01",
						Description: "desc01",
						Tags:        []string{"t01"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.AnythingOfType("transactionFunc")).
					Return(nil, nil)
			},
		},
		{
			name: "Delete exam02",
			args: &args{
				examId: examId02,
				userId: userId,
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id02,
						Topic:       "topic02",
						Description: "desc02",
						Tags:        []string{"t02"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.AnythingOfType("transactionFunc")).
					Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			err := s.examService.DeleteExam(
				args.examId,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestCreateQuestion() {
	type args struct {
		examId  string
		ask     string
		answers []string
		userId  string
	}

	type result struct {
		questionId string
		err        error
	}

	id01 := primitive.NewObjectID()
	examId01 := id01.Hex()

	id02 := primitive.NewObjectID()
	examId02 := id02.Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Create question01",
			args: &args{
				examId:  examId01,
				ask:     "q01",
				answers: []string{"a01"},
				userId:  "user01",
			},
			expected: &result{
				questionId: "question01",
				err:        nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id01,
						Topic:       "topic01",
						Description: "desc01",
						Tags:        []string{"t01"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CreateQuestion(mock.Anything, model.Question{
						ExamId:  args.examId,
						Ask:     args.ask,
						Answers: args.answers,
						UserId:  args.userId,
					}).
					Return("question01", nil)
			},
		},
		{
			name: "Create question02",
			args: &args{
				examId:  examId02,
				ask:     "q02",
				answers: []string{"a02"},
				userId:  "user01",
			},
			expected: &result{
				questionId: "question02",
				err:        nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id02,
						Topic:       "topic02",
						Description: "desc02",
						Tags:        []string{"t02"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CreateQuestion(mock.Anything, model.Question{
						ExamId:  args.examId,
						Ask:     args.ask,
						Answers: args.answers,
						UserId:  args.userId,
					}).
					Return("question02", nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			questionId, err := s.examService.CreateQuestion(
				args.examId,
				args.ask,
				args.answers,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.questionId, questionId)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestUpdaetQuestion() {
	type args struct {
		questionId string
		ask        string
		answers    []string
		userId     string
	}

	type result struct {
		questionId string
		err        error
	}

	id01 := primitive.NewObjectID()
	questionId01 := id01.Hex()

	id02 := primitive.NewObjectID()
	questionId02 := id02.Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Update question01 ask",
			args: &args{
				questionId: questionId01,
				ask:        "q01-updated",
				answers:    []string{"a01"},
				userId:     "user01",
			},
			expected: &result{
				questionId: questionId01,
				err:        nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetQuestionById(mock.Anything, args.questionId).
					Return(&model.Question{
						Id:      id01,
						ExamId:  "exam01",
						Ask:     "q01",
						Answers: []string{"a01"},
						UserId:  args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.AnythingOfType("transactionFunc")).
					Return(nil, nil)
			},
		},
		{
			name: "Update question02 answers",
			args: &args{
				questionId: questionId02,
				ask:        "q02",
				answers:    []string{"a01", "a02"},
				userId:     "user01",
			},
			expected: &result{
				questionId: questionId02,
				err:        nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetQuestionById(mock.Anything, args.questionId).
					Return(&model.Question{
						Id:      id02,
						ExamId:  "exam01",
						Ask:     "q02",
						Answers: []string{"a01"},
						UserId:  args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.AnythingOfType("transactionFunc")).
					Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			questionId, err := s.examService.UpdateQuestion(
				args.questionId,
				args.ask,
				args.answers,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.questionId, questionId)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindQuestions() {
	type args struct {
		pageIndex int32
		pageSize  int32
		examId    string
		userId    string
	}

	type result struct {
		total     int32
		pageCount int32
		questions []model.Question
		err       error
	}

	id := primitive.NewObjectID()
	examId := id.Hex()
	userId := "user01"
	mockQuestions := []model.Question{
		{},
		{},
		{},
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find questions by page1",
			args: &args{
				pageIndex: 0,
				pageSize:  10,
				examId:    examId,
				userId:    userId,
			},
			expected: &result{
				total:     3,
				pageCount: 1,
				questions: mockQuestions,
				err:       nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id,
						Topic:       "topic01",
						Description: "desc01",
						Tags:        []string{"t01"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					FindQuestionsByExamIdOrderByUpdateAtDesc(
						mock.Anything,
						args.examId,
						args.pageIndex*args.pageSize,
						args.pageSize).
					Return(mockQuestions, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, args.examId).
					Return(int32(3), nil)
			},
		},
		{
			name: "Find questions by page2",
			args: &args{
				pageIndex: 1,
				pageSize:  10,
				examId:    examId,
				userId:    userId,
			},
			expected: &result{
				total:     13,
				pageCount: 2,
				questions: mockQuestions,
				err:       nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id,
						Topic:       "topic01",
						Description: "desc01",
						Tags:        []string{"t01"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					FindQuestionsByExamIdOrderByUpdateAtDesc(
						mock.Anything,
						args.examId,
						args.pageIndex*args.pageSize,
						args.pageSize).
					Return(mockQuestions, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, args.examId).
					Return(int32(13), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			total, pageCount, questions, err := s.examService.FindQuestions(
				args.pageIndex,
				args.pageSize,
				args.examId,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.total, total)
			s.Equal(expected.pageCount, pageCount)
			s.Equal(expected.questions, questions)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestDeleteQuestion() {
	type args struct {
		questionId string
		userId     string
	}

	type result struct {
		err error
	}

	userId := "user01"
	id01 := primitive.NewObjectID()
	questionId01 := id01.Hex()

	id02 := primitive.NewObjectID()
	questionId02 := id02.Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Delete question01",
			args: &args{
				questionId: questionId01,
				userId:     userId,
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetQuestionById(mock.Anything, args.questionId).
					Return(&model.Question{
						Id:      id01,
						ExamId:  "exam01",
						Ask:     "ask01",
						Answers: []string{"a01", "a02"},
						UserId:  args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.AnythingOfType("transactionFunc")).
					Return(nil, nil)
			},
		},
		{
			name: "Delete question02",
			args: &args{
				questionId: questionId02,
				userId:     userId,
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetQuestionById(mock.Anything, args.questionId).
					Return(&model.Question{
						Id:      id02,
						ExamId:  "exam01",
						Ask:     "ask02",
						Answers: []string{"b01", "b02"},
						UserId:  args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.AnythingOfType("transactionFunc")).
					Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			err := s.examService.DeleteQuestion(args.questionId, args.userId)

			expected := tc.expected
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindRandomQuestions() {
	type args struct {
		examId string
		userId string
		size   int32
	}

	type result struct {
		exam      *model.Exam
		questions []model.Question
		err       error
	}

	id := primitive.NewObjectID()
	examId := id.Hex()
	userId := "user01"
	mockExam := model.Exam{
		Id:          id,
		Topic:       "topic01",
		Description: "desc01",
		Tags:        []string{},
		IsPublic:    true,
		UserId:      userId,
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find random questions by size 5",
			args: &args{
				examId: examId,
				userId: userId,
				size:   5,
			},
			expected: &result{
				exam:      &mockExam,
				questions: []model.Question{{}, {}, {}, {}, {}},
				err:       nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&mockExam, nil)
				s.mockDatabaseRepository.EXPECT().
					FindQuestionsByExamIdOrderByUpdateAtDesc(
						mock.Anything,
						args.examId,
						mock.AnythingOfType("int32"),
						int32(1)).
					Return([]model.Question{
						{},
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, examId).
					Return(int32(100), nil)
			},
		},
		{
			name: "Find random questions by size 10",
			args: &args{
				examId: examId,
				userId: userId,
				size:   10,
			},
			expected: &result{
				exam: &mockExam,
				questions: []model.Question{
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				},
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&mockExam, nil)
				s.mockDatabaseRepository.EXPECT().
					FindQuestionsByExamIdOrderByUpdateAtDesc(
						mock.Anything,
						args.examId,
						mock.AnythingOfType("int32"),
						int32(1)).
					Return([]model.Question{
						{},
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, examId).
					Return(int32(100), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			exam, questions, err := s.examService.FindRandomQuestions(
				args.examId,
				args.userId,
				args.size,
			)

			expected := tc.expected
			s.Equal(expected.exam, exam)
			s.Equal(expected.questions, questions)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestCreateExamRecord() {
	type args struct {
		examId           string
		score            int32
		wrongQuestionIds []string
		userId           string
	}

	type result struct {
		err error
	}

	userId := "user01"
	id := primitive.NewObjectID()
	examId := id.Hex()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Create examRecord01",
			args: &args{
				examId:           examId,
				score:            9,
				wrongQuestionIds: []string{"question01"},
				userId:           userId,
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id,
						Topic:       "topic",
						Description: "desc",
						Tags:        []string{"a01"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.Anything).
					Return(nil, nil)
			},
		},
		{
			name: "Create examRecord02",
			args: &args{
				examId:           examId,
				score:            8,
				wrongQuestionIds: []string{"question01", "question02"},
				userId:           userId,
			},
			expected: &result{
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(mock.Anything, args.examId).
					Return(&model.Exam{
						Id:          id,
						Topic:       "topic",
						Description: "desc",
						Tags:        []string{"a01"},
						IsPublic:    true,
						UserId:      args.userId,
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					WithTransaction(mock.Anything).
					Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			err := s.examService.CreateExamRecord(
				args.examId,
				args.score,
				args.wrongQuestionIds,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindExamRecords() {
	type args struct {
		pageIndex int32
		pageSize  int32
		examId    string
		userId    string
	}

	type result struct {
		total       int32
		pageCount   int32
		examRecords []model.ExamRecord
		err         error
	}

	id := primitive.NewObjectID()
	examId := id.Hex()
	userId := "user01"
	mockExamRecords := []model.ExamRecord{
		{},
		{},
		{},
	}

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find examRecords by page1",
			args: &args{
				pageIndex: 0,
				pageSize:  10,
				examId:    examId,
				userId:    userId,
			},
			expected: &result{
				total:       3,
				pageCount:   1,
				examRecords: mockExamRecords,
				err:         nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
						mock.Anything,
						args.examId,
						args.userId,
						args.pageIndex*args.pageSize,
						args.pageSize,
					).
					Return(mockExamRecords, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamRecordsByExamIdAndUserId(mock.Anything, args.examId, args.userId).
					Return(int32(3), nil)
			},
		},
		{
			name: "Find examRecords by page2",
			args: &args{
				pageIndex: 1,
				pageSize:  10,
				examId:    examId,
				userId:    userId,
			},
			expected: &result{
				total:       13,
				pageCount:   2,
				examRecords: mockExamRecords,
				err:         nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindExamRecordsByExamIdAndUserIdOrderByUpdateAtDesc(
						mock.Anything,
						args.examId,
						args.userId,
						args.pageIndex*args.pageSize,
						args.pageSize,
					).
					Return(mockExamRecords, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamRecordsByExamIdAndUserId(mock.Anything, args.examId, args.userId).
					Return(int32(13), nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			total, pageCount, examRecords, err := s.examService.FindExamRecords(
				args.pageIndex,
				args.pageSize,
				args.examId,
				args.userId,
			)

			expected := tc.expected
			s.Equal(expected.total, total)
			s.Equal(expected.pageCount, pageCount)
			s.Equal(expected.examRecords, examRecords)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindExamRecordOverview() {
	type args struct {
		examId    string
		userId    string
		startDate time.Time
	}

	type result struct {
		strStartDate string
		exam         *model.Exam
		questions    []model.Question
		answerWrongs []model.AnswerWrong
		examRecords  []model.ExamRecord
		err          error
	}

	userId := "user01"

	examId := "exam01"
	mockExam := &model.Exam{}
	mockQuestions := []model.Question{{}}
	mockWrongAnswers := []model.AnswerWrong{
		{
			ExamId:     examId,
			QuestionId: "q01",
			Times:      1,
			UserId:     userId,
		},
		{
			ExamId:     examId,
			QuestionId: "q02",
			Times:      2,
			UserId:     userId,
		},
	}
	mockExamRecords := []model.ExamRecord{{}}
	startDate := time.Now()

	examId02 := "exam02"
	mockExam02 := &model.Exam{}
	mockQuestions02 := []model.Question{}
	mockWrongAnswers02 := []model.AnswerWrong{}
	mockExamRecords02 := []model.ExamRecord{{}, {}}
	startDate02 := time.Now()

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find exam01 examRecordOverview",
			args: &args{
				examId:    examId,
				userId:    userId,
				startDate: startDate,
			},
			expected: &result{
				strStartDate: startDate.Format("2006/01/02"),
				exam:         mockExam,
				questions:    mockQuestions,
				answerWrongs: mockWrongAnswers,
				examRecords:  mockExamRecords,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(
						mock.Anything,
						args.examId,
					).
					Return(mockExam, nil)
				s.mockDatabaseRepository.EXPECT().
					FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc(
						mock.Anything,
						args.examId,
						args.userId,
						int32(10)).
					Return(mockWrongAnswers, nil)
				s.mockDatabaseRepository.EXPECT().
					FindQuestionsByQuestionIds(mock.Anything, []string{"q01", "q02"}).
					Return(mockQuestions, nil)
				s.mockDatabaseRepository.EXPECT().
					FindExamRecordsByExamIdAndUserIdAndCreatedAt(
						mock.Anything, args.examId, args.userId, args.startDate).
					Return(mockExamRecords, nil)
			},
		},
		{
			name: "Find exam02 examRecordOverview",
			args: &args{
				examId:    examId02,
				userId:    userId,
				startDate: startDate02,
			},
			expected: &result{
				strStartDate: startDate02.Format("2006/01/02"),
				exam:         mockExam02,
				questions:    mockQuestions02,
				answerWrongs: mockWrongAnswers02,
				examRecords:  mockExamRecords02,
				err:          nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					GetExamById(
						mock.Anything,
						args.examId,
					).
					Return(mockExam02, nil)
				s.mockDatabaseRepository.EXPECT().
					FindAnswerWrongsByExamIdAndUserIdOrderByTimesDesc(
						mock.Anything,
						args.examId,
						args.userId,
						int32(10)).
					Return(mockWrongAnswers02, nil)
				s.mockDatabaseRepository.EXPECT().
					FindQuestionsByQuestionIds(mock.Anything, []string{}).
					Return(mockQuestions02, nil)
				s.mockDatabaseRepository.EXPECT().
					FindExamRecordsByExamIdAndUserIdAndCreatedAt(
						mock.Anything, args.examId, args.userId, args.startDate).
					Return(mockExamRecords02, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			strStartDate, exam, questions, answerWrongs, examRecords, err := s.examService.FindExamRecordOverview(
				args.examId,
				args.userId,
				args.startDate,
			)

			expected := tc.expected
			s.Equal(expected.strStartDate, strStartDate)
			s.Equal(expected.exam, exam)
			s.Equal(expected.questions, questions)
			s.Equal(expected.answerWrongs, answerWrongs)
			s.Equal(expected.examRecords, examRecords)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestFindExamInfos() {
	type args struct {
		userId   string
		isPublic bool
	}

	type result struct {
		examInfos []ExamInfo
		err       error
	}

	userId := "user01"

	id01 := primitive.NewObjectID()
	examId01 := id01.Hex()

	id02 := primitive.NewObjectID()
	examId02 := id02.Hex()

	var (
		questionCount1 int32 = 10
		questionCount2 int32 = 20
		recordCount1   int32 = 1
		recordCount2   int32 = 2
	)

	testCases := []struct {
		name     string
		args     *args
		expected *result
		on       func(s *MyTestSuite, args *args)
	}{
		{
			name: "Find examInfos by isPublic true",
			args: &args{
				userId:   userId,
				isPublic: true,
			},
			expected: &result{
				examInfos: []ExamInfo{
					{
						ExamId:        examId01,
						Topic:         "topic01",
						Description:   "desc01",
						IsPublic:      true,
						QuestionCount: questionCount1,
						RecordCount:   recordCount1,
					},
					{
						ExamId:        examId02,
						Topic:         "topic02",
						Description:   "desc02",
						IsPublic:      true,
						QuestionCount: questionCount2,
						RecordCount:   recordCount2,
					},
				},
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(
						mock.Anything, args.userId, args.isPublic).
					Return([]model.Exam{
						{
							Id:          id01,
							Topic:       "topic01",
							Description: "desc01",
							Tags:        []string{"a01"},
							IsPublic:    args.isPublic,
							UserId:      args.userId,
						},
						{
							Id:          id02,
							Topic:       "topic02",
							Description: "desc02",
							Tags:        []string{"a02"},
							IsPublic:    args.isPublic,
							UserId:      args.userId,
						},
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, examId01).
					Return(questionCount1, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, examId02).
					Return(questionCount2, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamRecordsByExamIdAndUserId(mock.Anything, examId01, args.userId).
					Return(recordCount1, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamRecordsByExamIdAndUserId(mock.Anything, examId02, args.userId).
					Return(recordCount2, nil)
			},
		},
		{
			name: "Find examInfos by isPublic false",
			args: &args{
				userId:   userId,
				isPublic: false,
			},
			expected: &result{
				examInfos: []ExamInfo{
					{
						ExamId:        examId01,
						Topic:         "topic01",
						Description:   "desc01",
						IsPublic:      false,
						QuestionCount: questionCount1,
						RecordCount:   recordCount1,
					},
				},
				err: nil,
			},
			on: func(s *MyTestSuite, args *args) {
				s.mockDatabaseRepository.EXPECT().
					FindExamsByUserIdAndIsPublicOrderByUpdateAtDesc(
						mock.Anything, args.userId, args.isPublic).
					Return([]model.Exam{
						{
							Id:          id01,
							Topic:       "topic01",
							Description: "desc01",
							Tags:        []string{"a01"},
							IsPublic:    args.isPublic,
							UserId:      args.userId,
						},
					}, nil)
				s.mockDatabaseRepository.EXPECT().
					CountQuestionsByExamId(mock.Anything, examId01).
					Return(questionCount1, nil)
				s.mockDatabaseRepository.EXPECT().
					CountExamRecordsByExamIdAndUserId(mock.Anything, examId01, args.userId).
					Return(recordCount1, nil)
			},
		},
	}

	for _, tc := range testCases {
		s.SetupTest()
		s.Run(tc.name, func() {
			args := tc.args
			tc.on(s, args)

			// Test
			examInfos, err := s.examService.FindExamInfos(
				args.userId,
				args.isPublic,
			)

			expected := tc.expected
			s.Equal(expected.examInfos, examInfos)
			s.Equal(expected.err, err)
		})
	}
}

func (s *MyTestSuite) TestCircleCI() {
	// test 6
	s.Fail("========= For test CircleCI !!!")
}
