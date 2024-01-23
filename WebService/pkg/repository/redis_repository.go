package repository

import (
	context "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kakurineuin/learn-english-microservices/web-service/pb"
)

const KEY_WORD_MEANING = "word-meaning:"

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository() *RedisRepository {
	return &RedisRepository{}
}

func (repo *RedisRepository) ConnectDB(uri string) error {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return fmt.Errorf("ConnectDB failed! error: %w", err)
	}

	rdb := redis.NewClient(opt)
	repo.client = rdb

	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := repo.client.Ping(pingCtx).Err(); err != nil {
		return fmt.Errorf("ConnectDB ping Redis failed! error: %w", err)
	}

	fmt.Println("Connected to Redis")
	return nil
}

func (repo *RedisRepository) DisconnectDB() error {
	if err := repo.client.Close(); err != nil {
		return fmt.Errorf("DisconnectDB failed! error: %w", err)
	}

	return nil
}

func (repo *RedisRepository) CreateWordMeanings(
	ctx context.Context,
	key string,
	wordMeanings []*pb.WordMeaning,
	expiration time.Duration,
) error {
	bytes, err := json.Marshal(wordMeanings)
	if err != nil {
		return err
	}

	fullKey := KEY_WORD_MEANING + key
	_, err = repo.client.Set(ctx, fullKey, bytes, expiration).Result()
	return err
}

func (repo *RedisRepository) FindWordMeanings(
	ctx context.Context,
	key string,
) (wordMeanings []*pb.WordMeaning, err error) {
	data, err := repo.client.Get(ctx, KEY_WORD_MEANING+key).Result()
	if err != nil {

		// Key 不存在
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &wordMeanings)
	if err != nil {
		return nil, err
	}

	return wordMeanings, nil
}
