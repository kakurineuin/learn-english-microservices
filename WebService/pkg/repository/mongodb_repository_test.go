package repository

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kakurineuin/learn-english-microservices/web-service/pkg/model"
)

// 使用測試的資料庫
const DATABASE = "learnEnglish_test"

type MyTestSuite struct {
	suite.Suite
	repo             DatabaseRepository
	uri              string
	mongodbContainer *mongodb.MongoDBContainer
	client           *mongo.Client
	userCollection   *mongo.Collection
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
	err = s.repo.ConnectDB(ctx, uri)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.uri = uri
	s.mongodbContainer = mongodbContainer

	// 用來建立測試資料的 client
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri).SetTimeout(10*time.Second),
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.client = client
	s.userCollection = client.Database(DATABASE).Collection("users")
}

// run once, after test suite methods
func (s *MyTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
	ctx := context.Background()

	// 不呼叫 panic，為了繼續往下執行去關閉 container
	if err := s.client.Disconnect(ctx); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// 不呼叫 panic，為了繼續往下執行去關閉 container
	if err := s.repo.DisconnectDB(ctx); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// Terminate container
	if err := s.mongodbContainer.Terminate(ctx); err != nil {
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
	ctx := context.Background()
	repo := NewMongoDBRepository(DATABASE)

	err := repo.ConnectDB(ctx, s.uri)
	s.Nil(err)

	err = repo.DisconnectDB(ctx)
	s.Nil(err)
}

func (s *MyTestSuite) TestCreateUser() {
	ctx := context.Background()

	userId, err := s.repo.CreateUser(ctx, model.User{
		Username: "TestCreateUser",
		Password: "TestCreateUser",
		Role:     "user",
	})
	s.Nil(err)
	s.NotEmpty(userId)
}

func (s *MyTestSuite) TestGetUserById() {
	ctx := context.Background()

	result, err := s.userCollection.InsertOne(ctx, model.User{
		Username: "TestCreateUser",
		Password: "TestCreateUser",
		Role:     "user",
	})
	s.Nil(err)
	userId := result.InsertedID.(primitive.ObjectID).Hex()

	user, err := s.repo.GetUserById(ctx, userId)
	s.Nil(err)
	s.NotNil(user)
}

func (s *MyTestSuite) TestGetUserByUsername() {
	ctx := context.Background()

	username := "TestCreateUser"

	_, err := s.userCollection.InsertOne(ctx, model.User{
		Username: username,
		Password: "TestCreateUser",
		Role:     "user",
	})
	s.Nil(err)

	user, err := s.repo.GetUserByUsername(ctx, username)
	s.Nil(err)
	s.NotNil(user)
}

func (s *MyTestSuite) TestGetAdminUser() {
	ctx := context.Background()

	_, err := s.userCollection.InsertOne(ctx, model.User{
		Username: "TestGetAdminUser",
		Password: "TestGetAdminUser",
		Role:     "admin",
	})
	s.Nil(err)

	user, err := s.repo.GetAdminUser(ctx)
	s.Nil(err)
	s.NotNil(user)
}
