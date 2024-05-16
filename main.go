package main

import (
	"fmt"
	"net/http"
	"crypto/sha256"
    "encoding/base64"
    "net/url"
    "strings"

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
	fmt.Printf("On urlSubmission " + fmt.Sprintf("%s", url))
	shortURL, err := processURL(url, svc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(fmt.Sprintf("%s", shortURL))

	c.JSON(http.StatusOK, gin.H{"message": "URL stored successfully", "shortenedUrl": shortURL})
}

func generateBase64Path(originalURL string) string {
    // Parse the original URL
    parsedURL, err := url.Parse(originalURL)
    if err != nil {
        // handle error
        return ""
    }

    // Hash the URL path using SHA-256
    hasher := sha256.New()
    hasher.Write([]byte(parsedURL.Path))
    hash := hasher.Sum(nil)

    // Encode the hash using base64
    base64EncodedHash := base64.URLEncoding.EncodeToString(hash)

    // Ensure the encoded hash matches the regex pattern /[a-zA-Z0-9]{9}/
    base64EncodedHash = strings.TrimRight(base64EncodedHash, "=")
    if len(base64EncodedHash) > 9 {
        base64EncodedHash = base64EncodedHash[:9]
    } else if len(base64EncodedHash) < 9 {
        base64EncodedHash = strings.Repeat("a", 9-len(base64EncodedHash)) + base64EncodedHash
    }

    return base64EncodedHash
}

func processURL(longURL string, svc *dynamodb.DynamoDB) (string, error) {
	fmt.Printf("Processing %s\n", longURL)

	// Generate a shortened URL (you'll need to implement this logic)
	shortenedURLPath := generateBase64Path(longURL)

	item := map[string]*dynamodb.AttributeValue{
		"url":          {S: aws.String(longURL)},
		"shortenedUrl": {S: aws.String(shortenedURLPath)},
		// Add any other attributes as needed
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
    fmt.Printf("Loading short url %s\n", shortenedURL)

    // Query DynamoDB to get the original URL
    input := &dynamodb.GetItemInput{
        TableName: aws.String("UrlMap"),
        Key: map[string]*dynamodb.AttributeValue{
            "shortenedUrl": {S: aws.String(shortenedURL)},
        },
    }

    result, err := svc.GetItem(input)
    if err != nil {
        fmt.Println("Error retrieving item from DynamoDB:", err)
        c.String(http.StatusInternalServerError, "Error retrieving URL")
        return
    }

    if result.Item == nil {
        fmt.Println("URL not found in DynamoDB")
        c.String(http.StatusNotFound, "URL not found")
        return
    }

    // Extract the original URL from the DynamoDB response
    originalURL := *result.Item["url"].S
    fmt.Printf("Retrieved original URL: %s\n", originalURL)

    // Redirect the user to the original URL
    c.Redirect(http.StatusMovedPermanently, originalURL)
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