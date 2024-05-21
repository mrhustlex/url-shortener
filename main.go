package main

import (
	"crypto/sha1"
    "fmt"
    "math/rand"
    "time"
	"net/http"
		"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)
var svc *dynamodb.DynamoDB

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
	// r.GET("/redirect", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "http://www.aws.amazon.com")
	// })
	// Shorten redirect URL
	r.GET("/:url", urlRetrieval)

	return r
}

func urlSubmission(c *gin.Context) {
	// Parse the JSON request body
	url := c.Request.FormValue("url")
	fmt.Printf("On urlSubmission " + fmt.Sprintf("%s", url))
	shortURL, err := processURL(url, svc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(fmt.Sprintf("%s", shortURL))

	c.JSON(http.StatusOK, gin.H{"message": "URL stored successfully", "shortenedUrl": shortURL})
}


var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateShortPath(longURL string) string {
    // Generate a unique string from the long URL using SHA-1 hashing
    h := sha1.New()
    h.Write([]byte(longURL))
    hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

    // Use current time to add randomness to the string
    rand.Seed(time.Now().UnixNano())

    // Generate a random string of length 9 using the letterRunes
    b := make([]rune, 9)
    for i := range b {
        b[i] = letterRunes[rand.Intn(len(letterRunes))]
    }
    randomString := string(b)

    // Combine the first 8 characters of the hash and the random string to create the short URL
    shortURL := hash[:8] + randomString[:1]

    return shortURL
}

func processURL(longURL string, svc *dynamodb.DynamoDB) (string, error) {
	fmt.Printf("Processing %s\n", longURL)
	// Generate a shortened URL (you'll need to implement this logic)
	shortenedURLPath := generateShortPath(longURL)

	fmt.Printf("Short path: %s\n", shortenedURLPath)

	item := map[string]*dynamodb.AttributeValue{
		"url":          {S: aws.String(longURL)},
		"shortenedUrl": {S: aws.String(shortenedURLPath)},
		"Count":		{N: aws.String("1")},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("UrlMap"),
		Item:      item,
	}

	_, err := svc.PutItem(input)
	if err != nil {
		return "", fmt.Errorf("error putting item in DynamoDB: %w", err)
	}

	fmt.Println("Item put successfully!")
	return shortenedURLPath, nil
}




func urlRetrieval(c *gin.Context) {
    fmt.Println("On urlRetrieval")

    // Get the shortened URL from the request path parameter
    shortenedURL := c.Param("url")
	params := &dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{ 
            "shortenedUrl": { 
                S: aws.String(shortenedURL),
            },
        },
        TableName: aws.String("UrlMap"), 
        ConsistentRead: aws.Bool(true),
    }
    resp, err := svc.GetItem(params)

    if err != nil {
        fmt.Println("Query failed", err)
    }else{
		originalURL := *resp.Item["url"].S
		fmt.Printf("Retrieved original URL: %s\n", originalURL)
		// // Redirect the user to the original URL
		c.Redirect(http.StatusMovedPermanently, "http://www.google.com")
	}

	
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
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		for _, n := range result.TableNames {
			fmt.Println(*n)
			c.String(http.StatusOK, *n)
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