package ristretto

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/ristretto"

	"github.com/grokify/sogo/database/kvs"
)

type Client struct {
	cache *ristretto.Cache
}

func NewClient(cfg kvs.Config) (*Client, error) {
	// Set default cache configuration if not provided
	config := &ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	}

	cache, err := ristretto.NewCache(config)
	if err != nil {
		return nil, err
	}

	return &Client{cache: cache}, nil
}

func (c *Client) SetString(ctx context.Context, key, val string) error {
	// Context is not used by ristretto as it's an in-memory cache
	// but we accept it to conform to the interface
	ok := c.cache.Set(key, val, 1)
	if !ok {
		return errors.New("failed to set value in cache")
	}
	return nil
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	// Context is not used by ristretto as it's an in-memory cache
	// but we accept it to conform to the interface
	val, found := c.cache.Get(key)
	if !found {
		return "", errors.New("key not found")
	}

	str, ok := val.(string)
	if !ok {
		return "", errors.New("value is not a string")
	}

	return str, nil
}

func (c *Client) GetOrDefaultString(ctx context.Context, key, def string) string {
	val, err := c.GetString(ctx, key)
	if err != nil {
		return def
	}
	return val
}

func (c *Client) SetAny(ctx context.Context, key string, val any) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return c.SetString(ctx, key, string(bytes))
}

func (c *Client) GetAny(ctx context.Context, key string, val any) error {
	str, err := c.GetString(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), val)
}

func (c *Client) Close() {
	c.cache.Close()
}
