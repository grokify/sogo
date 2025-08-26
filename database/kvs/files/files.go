package files

import (
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

func (client Client) SetString(key, val string) error {
	// tempval, err2 = strconv.ParseUint(data["Perm"], 10, 32)
	return os.WriteFile(
		KeyToFilename(key),
		[]byte(val),
		os.FileMode(client.config.FileMode))
}

func (client Client) GetString(key string) (string, error) {
	data, err := os.ReadFile(KeyToFilename(key))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (client Client) GetOrEmptyString(key string) string {
	val, err := client.GetString(key)
	if err != nil {
		return ""
	}
	return val
}
