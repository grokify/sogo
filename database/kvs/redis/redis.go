package redis

import (
	"context"
	"encoding/json"

	redis "github.com/go-redis/redis/v8"
	"github.com/grokify/sogo/database/kvs"
)

type Client struct {
	redisClient *redis.Client
}

func NewClient(cfg kvs.Config) *Client {
	return &Client{
		redisClient: redis.NewClient(NewRedisOptions(cfg))}
}

func (client Client) SetString(key, val string) error {
	// For context, see https://github.com/go-redis/redis/issues/582
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	return client.redisClient.Set(context.Background(), key, val, 0).Err()
}

func (client Client) GetString(key string) (string, error) {
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	return client.redisClient.Get(context.Background(), key).Result()
}

func (client Client) GetOrEmptyString(key string) string {
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	if val, err := client.redisClient.Get(context.Background(), key).Result(); err != nil {
		return ""
	} else {
		return val
	}
}

func (client Client) SetInterface(key string, val interface{}) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return client.redisClient.Set(
		context.Background(), key, string(bytes), 0).Err()
}

func (client Client) GetInterface(key string, val interface{}) error {
	strCmd := client.redisClient.Get(context.Background(), key)
	bytes, err := strCmd.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, val)
}

func NewRedisOptions(cfg kvs.Config) *redis.Options {
	return &redis.Options{
		Addr:     cfg.HostPort(),
		Password: cfg.Password,
		DB:       cfg.CustomIndex}
}
