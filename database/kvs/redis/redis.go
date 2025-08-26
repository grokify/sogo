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

func (client Client) SetString(ctx context.Context, key, val string) error {
	// For context, see https://github.com/go-redis/redis/issues/582
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	return client.redisClient.Set(ctx, key, val, 0).Err()
}

func (client Client) GetString(ctx context.Context, key string) (string, error) {
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	return client.redisClient.Get(ctx, key).Result()
}

func (client Client) GetOrDefaultString(ctx context.Context, key, def string) string {
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	if val, err := client.redisClient.Get(ctx, key).Result(); err != nil {
		return def
	} else {
		return val
	}
}

func (client Client) SetAny(ctx context.Context, key string, val any) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return client.redisClient.Set(ctx, key, string(bytes), 0).Err()
}

func (client Client) GetAny(ctx context.Context, key string, val any) error {
	strCmd := client.redisClient.Get(ctx, key)
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
