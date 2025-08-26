package dynamodb

// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-create-table-item.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/read_item.go

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/grokify/sogo/database/kvs"
)

const (
	KeyName   = "key"
	ValueName = "value"
)

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Client struct {
	config         kvs.Config
	dynamodbClient *dynamodb.Client
}

func NewClient(cfg kvs.Config) (*Client, error) {
	cfg.Region = strings.TrimSpace(cfg.Region)
	if len(cfg.Region) == 0 {
		return nil, errors.New("E_NO_REGION_FOR_AWS")
	}
	cfg.Table = strings.TrimSpace(cfg.Table)
	if len(cfg.Table) == 0 {
		return nil, errors.New("E_NO_TABLE_FOR_DYNAMODB")
	}
	if cfg.DynamodbReadUnits == 0 {
		cfg.DynamodbReadUnits = 1
	}
	if cfg.DynamodbWriteUnits == 0 {
		cfg.DynamodbWriteUnits = 1
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.Region))
	if err != nil {
		return nil, err
	}

	return &Client{
		config:         cfg,
		dynamodbClient: dynamodb.NewFromConfig(awsCfg)}, nil
}

func (client Client) SetString(key, val string) error {
	item := Item{
		Key:   key,
		Value: val}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(client.config.Table)}

	_, err = client.dynamodbClient.PutItem(context.TODO(), input)
	return err
}

func (client Client) GetString(key string) (string, error) {
	result, err := client.dynamodbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(client.config.Table),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{
				Value: key,
			},
		},
	})
	if err != nil {
		return "", err
	}
	item := Item{}

	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}
	return item.Value, nil
}

func (client Client) GetOrEmptyString(key string) string {
	val, err := client.GetString(key)
	if err != nil {
		return ""
	}
	return val
}

func (client Client) CreateTable() (*dynamodb.CreateTableOutput, error) {
	return client.dynamodbClient.CreateTable(context.TODO(), client.createTableInput())
}

func (client Client) createTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(KeyName),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(ValueName),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(KeyName),
				KeyType:       types.KeyTypeHash,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(client.config.DynamodbReadUnits),
			WriteCapacityUnits: aws.Int64(client.config.DynamodbWriteUnits),
		},
		TableName: aws.String(client.config.Table),
	}
}

