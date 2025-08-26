package files

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/grokify/sogo/database/kvs"
)

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Client struct {
	config kvs.Config
}

func NewClient(cfg kvs.Config) *Client {
	return &Client{config: cfg}
}

func KeyToFilename(key string) string {
	return strings.TrimSpace(key) + ".txt"
}

func (client Client) SetString(ctx context.Context, key, val string) error {
	// tempval, err2 = strconv.ParseUint(data["Perm"], 10, 32)
	return os.WriteFile(
		KeyToFilename(key),
		[]byte(val),
		os.FileMode(client.config.FileMode))
}

func (client Client) GetString(ctx context.Context, key string) (string, error) {
	data, err := os.ReadFile(KeyToFilename(key))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (client Client) GetOrDefaultString(ctx context.Context, key, def string) string {
	val, err := client.GetString(ctx, key)
	if err != nil {
		return def
	}
	return val
}

func (client Client) SetAny(ctx context.Context, key string, val any) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return client.SetString(ctx, key, string(bytes))
}

func (client Client) GetAny(ctx context.Context, key string, val any) error {
	str, err := client.GetString(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), val)
}
