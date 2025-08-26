package dynamodb

// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-create-table-item.html
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html
// https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/read_item.go

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
	dynamodbClient *dynamodb.DynamoDB
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

	sess, err := session.NewSession(NewAwsConfig(cfg))
	if err != nil {
		return nil, err
	}

	return &Client{
		config:         cfg,
		dynamodbClient: dynamodb.New(sess)}, nil
}

func (client Client) SetString(key, val string) error {
	item := Item{
		Key:   key,
		Value: val}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(client.config.Table)}

	_, err = client.dynamodbClient.PutItem(input)
	return err
}

func (client Client) GetString(key string) (string, error) {
	result, err := client.dynamodbClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(client.config.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return "", err
	}
	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
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
	return client.dynamodbClient.CreateTable(client.createTableInput())
}

func (client Client) createTableInput() *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(KeyName),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String(ValueName),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(KeyName),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(client.config.DynamodbReadUnits),
			WriteCapacityUnits: aws.Int64(client.config.DynamodbWriteUnits),
		},
		TableName: aws.String(client.config.Table),
	}
}

func NewAwsConfig(cfg kvs.Config) *aws.Config {
	return &aws.Config{
		Region: aws.String(cfg.Region)}
}
