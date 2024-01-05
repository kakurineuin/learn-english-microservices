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
	ctx              context.Context
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
	s.userCollection = client.Database(DATABASE).Collection("users")
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

func (s *MyTestSuite) TestCreateUser() {
	ctx := context.TODO()

	userId, err := s.repo.CreateUser(ctx, model.User{
		Username: "TestCreateUser",
		Password: "TestCreateUser",
		Role:     "user",
	})
	s.Nil(err)
	s.NotEmpty(userId)
}

func (s *MyTestSuite) TestGetUserById() {
	ctx := context.TODO()

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
	ctx := context.TODO()

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
	ctx := context.TODO()

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
