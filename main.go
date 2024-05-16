package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})
	// URL submission
	r.POST("/new_url", urlSubmission)
	// get all tables
	r.GET("/getTables", func(c *gin.Context) {
		getAllTables(c, svc)
	})
	// Shorten redirect URL
	r.GET("/:url", urlRetrieval)


	return r
}

func urlSubmission(c *gin.Context) {
	// Parse the JSON request body
	url := c.Request.FormValue("url")
	fmt.Printf("Processing " + fmt.Sprintf("%s", url))

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No URL input provided"})
		return
	}

	// Create a new item to be inserted into the DynamoDB table
	item := map[string]*dynamodb.AttributeValue{
		"url":             {S: aws.String(url)},
		"shorternedUrl": {S: aws.String(generateShortURL(url))}, // Replace generateShortURL with your own logic for generating short URLs
	}

	// Create the input configuration instance
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("UrlMap"), // Replace with your actual table name
	}

	// Put the item into the DynamoDB table
	_, err := svc.PutItem(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL stored successfully"})
}

func generateShortURL(longURL string) string {
    // Generate a short URL based on the long URL
    // Implement your logic here

    shortURL := "http://shorturl.com/abcd" // Replace with your actual implementation

    return shortURL
}

func urlRetrieval(c *gin.Context) {
	fmt.Println("On urlRetrieval")
	// c.Redirect(http.StatusFound, "http://www.google.com/")
	c.Redirect(http.StatusFound, "/new-url")
}

func getAllTables(c *gin.Context, svc *dynamodb.DynamoDB) {
	// Create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	fmt.Printf("Tables:\n")

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
		}

		// Assign the last read tablename as the start for our next call to the ListTables function
		// The maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
}

var svc *dynamodb.DynamoDB

func main() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000")})

	r := setupRouter()
	r.Run(":8080")
}