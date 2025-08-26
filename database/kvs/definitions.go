package kvs

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Config struct {
	Host               string `json:"host,omitempty"`
	Port               int    `json:"port,omitempty"`
	Password           string `json:"password,omitempty"`
	CustomIndex        int    `json:"customIndex,omitempty"` // Redis Database Index
	Region             string `json:"region,omitempty"`      // DynamoDB
	Table              string `json:"table,omitempty"`       // DynamoDB
	Directory          string `json:"directory,omitempty"`
	FileMode           uint32 `json:"fileMode,omitempty"`
	DynamodbReadUnits  int64  `json:"dynamodbReadUnits,omitempty"`
	DynamodbWriteUnits int64  `json:"dynamodbWriteUnits,omitempty"`
}

func ParseConfig(jsonBytes []byte) (*Config, error) {
	var cfg Config
	err := json.Unmarshal(jsonBytes, &cfg)
	return &cfg, err
}

func (cfg *Config) HostPort() string {
	parts := []string{}
	cfg.Host = strings.TrimSpace(cfg.Host)
	if len(cfg.Host) > 0 {
		parts = append(parts, cfg.Host)
	}
	if (cfg.Port) > 0 {
		parts = append(parts, strconv.Itoa(cfg.Port))
	}
	return strings.Join(parts, ":")
}

type Client interface {
	SetString(key, val string) error
	GetString(key string) (string, error)
	GetOrEmptyString(key string) string
	SetInterface(key string, val interface{}) error
	GetInterface(key string, val interface{}) error
}
