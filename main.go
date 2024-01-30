package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	debugMode          bool
	dynamoDBTableName  string
	listeningPort      string
	defaultDynamoDBTable = "argocd"
	logger              *zap.Logger
)

// Request represents the JSON structure {"input": {"parameters": {"ddb-table-name": "name"}}}
type RequestBody struct {
	Input Input `json:"input"`
}

// Input represents the JSON structure {"parameters": {"ddb-table-name": "name"}}
type Input struct {
	Parameters Parameters `json:"parameters"`
}

// Parameters represents the JSON structure {"ddb-table-name": "name"}
type Parameters struct {
	DDBTableName string `json:"ddb-table-name"`
}

func init() {
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.StringVar(&dynamoDBTableName, "ddb-table-name", defaultDynamoDBTable, "DynamoDB table name (default: argocd)")
	flag.StringVar(&listeningPort, "listening-port", "8080", "Listening port for the web server")

	// Initialize logger
	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func main() {
	flag.Parse()

	// Print flag values in debug mode
	if debugMode {
		logger.Info("Debug mode is enabled",
			zap.String("DynamoDBTableName", dynamoDBTableName),
			zap.String("ListeningPort", listeningPort),
		)
	}

	// Check if required environment variables are set
	if dynamoDBTableName == "" {
		logger.Error("ddb-table-name environment variable is required")
		os.Exit(1)
	}

	// // Get the AWS region from environment variable
	awsRegion := os.Getenv("AWS_DEFAULT_REGION")
	// if awsRegion == "" {
	// 	logger.Error("AWS_DEFAULT_REGION environment variable is required")
	// 	os.Exit(1)
	// }
	
	if awsRegion == "" {
		// Use AWS SDK to determine the region
		sess, err := session.NewSession(&aws.Config{})
		if err != nil {
			logger.Fatal("Failed to create AWS session", zap.Error(err))
		}

		region := aws.StringValue(sess.Config.Region)
		os.Setenv("AWS_DEFAULT_REGION", region)
		logger.Info("AWS_DEFAULT_REGION environment variable set", zap.String("Region", region))
	
	}

	// Get the listening port from environment variable if present
	envListeningPort := os.Getenv("LISTENING_PORT")
	if envListeningPort != "" {
		listeningPort = envListeningPort

	}

	// AWS session creation with the specified region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		logger.Fatal("Failed to create AWS session", zap.Error(err))
	}

	// Set debug mode
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// DynamoDB client creation
	dynamoDBClient := dynamodb.New(sess)

	// Setup Gin router
	router := gin.Default()

	// Define route for "/api/v1/getparams.execute"
	router.POST("/api/v1/getparams.execute", func(c *gin.Context) {
		// Parse request body
		var requestBody RequestBody
		if err := c.BindJSON(&requestBody); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body"})
			return
		}
		logger.Info("got request", zap.Any("request",requestBody ))
		// Use request body parameters if present
		if requestBody.Input.Parameters.DDBTableName != "" {
			dynamoDBTableName = requestBody.Input.Parameters.DDBTableName
		}

		// Retrieve data from DynamoDB
		params, err := getDynamoDBData(dynamoDBClient, dynamoDBTableName)
		if err != nil {
			logger.Error("Failed to retrieve data from DynamoDB", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data from DynamoDB"})
			return
		}

		// Respond with the retrieved data
		c.JSON(http.StatusOK, gin.H{"output": gin.H{"parameters": params}})
	})

	// Run the web server
	err = router.Run(":" + listeningPort)
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func getDynamoDBData(client *dynamodb.DynamoDB, tableName string) ([]interface{}, error) {
// func getDynamoDBData(client *dynamodb.DynamoDB, tableName string) ([]map[string]*dynamodb.AttributeValue, error) {
	
	// DynamoDB query input
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Execute the query
	result, err := client.Scan(input)
	if err != nil {
		logger.Error("Failed to scan DynamoDB table", zap.Error(err))
		return nil, err
	}

	// Convert DynamoDB items to map[string]interface{}
	// var items []map[string]*dynamodb.AttributeValue

	// items := make([]map[string]interface{}, len(result.Items))
	items := make([]interface{}, len(result.Items))
	for i, item := range result.Items {
		attributeValues := make(map[string]*dynamodb.AttributeValue)
		for key, value := range item {
			attributeValues[key] = value
		}

		var itemMap map[string]interface{}
		if err := dynamodbattribute.UnmarshalMap(attributeValues, &itemMap); err != nil {
			logger.Error("Failed to unmarshal DynamoDB attributes", zap.Error(err))
			return nil, err
		}

		items[i] = itemMap
	}

	logger.Info("output request", zap.Any("items", items))
	// return map[string]interface{}{"items": items}, nil
	return items, nil

}
