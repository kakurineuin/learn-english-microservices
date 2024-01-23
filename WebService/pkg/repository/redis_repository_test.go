package repository

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
)

type RedisRepositoryTestSuite struct {
	suite.Suite
	repo           CacheRepository
	uri            string
	redisContainer *tcredis.RedisContainer
	client         *redis.Client
}

func TestRedisRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RedisRepositoryTestSuite))
}

// run once, before test suite methods
func (s *RedisRepositoryTestSuite) SetupSuite() {
	log.Println("SetupSuite()")

	// Run container
	ctx := context.Background()
	redisContainer, err := tcredis.RunContainer(ctx, testcontainers.WithImage("docker.io/redis:7"))
	if err != nil {
		s.FailNow(err.Error())
	}

	uri, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.repo = NewRedisRepository()
	err = s.repo.ConnectDB(uri)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.uri = uri
	s.redisContainer = redisContainer

	// 用來建立測試資料的 client
	opt, err := redis.ParseURL(uri)
	if err != nil {
		s.FailNow(err.Error())
	}

	client := redis.NewClient(opt)
	s.client = client
}

// run once, after test suite methods
func (s *RedisRepositoryTestSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
	ctx := context.Background()

	// 不呼叫 panic，為了繼續往下執行去關閉 container
	if err := s.client.Close(); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// 不呼叫 panic，為了繼續往下執行去關閉 container
	if err := s.repo.DisconnectDB(); err != nil {
		log.Printf("DisconnectDB error: %v", err)
	}

	// Terminate container
	if err := s.redisContainer.Terminate(ctx); err != nil {
		log.Printf("redisContainer.Terminate() error: %v", err)
	}
}

// run before each test
func (s *RedisRepositoryTestSuite) SetupTest() {
	log.Println("SetupTest()")
}

// run after each test
func (s *RedisRepositoryTestSuite) TearDownTest() {
	log.Println("TearDownTest()")
}

// run before each test
func (s *RedisRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
}

// run after each test
func (s *RedisRepositoryTestSuite) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}

func (s *RedisRepositoryTestSuite) TestConnectDBAndDisconnectDB() {
	repo := NewRedisRepository()

	err := repo.ConnectDB(s.uri)
	s.Nil(err)

	err = repo.DisconnectDB()
	s.Nil(err)
}

func (s *RedisRepositoryTestSuite) TestCreateWordMeanings() {
	ctx := context.Background()

	wordMeanings := []*pb.WordMeaning{
		{
			Word: "test",
		},
	}
	err := s.repo.CreateWordMeanings(ctx, "key01", wordMeanings, 5*time.Minute)
	s.Nil(err)
}

func (s *RedisRepositoryTestSuite) TestFindWordMeanings() {
	ctx := context.Background()

	bytes, err := json.Marshal([]*pb.WordMeaning{
		{
			Word: "test",
		},
		{
			Word: "test",
		},
	})
	s.Nil(err)

	key := "key01"
	_, err = s.client.Set(ctx, KEY_WORD_MEANING+key, bytes, 5*time.Minute).Result()
	s.Nil(err)

	wordMeanings, err := s.repo.FindWordMeanings(ctx, key)
	s.Nil(err)
	s.Len(wordMeanings, 2)
}
