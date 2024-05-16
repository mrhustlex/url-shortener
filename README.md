# url-shortener
This project is working on an URL shortener powered by GoLang
## System Design

### Functional Requirement
* Shorten an URL
* Retrieve an URL with the shortened URL
### Non-Functional Requirement
* High Availability
  * Automatically recovery when the service health check is failed
  * Eliminate single point of failure
* Scalability
  * Scale horizontally when the workload is getting more intense
  * Handle concurrent connection
* Security
  * DDoS protection
* CI/CD consideration

  
## Technology 
Programing language: GoLang with Gin framework
Server: Nginx
Cloud technology: Load Balancer
Database: NoSQL(Amazon DynamoDB)

## API design
POST/new_url
Assuming the NoSQL Database is a hashmap with original url(oUrl) as key and shortened URL(sURL) as value

#### Sudo code logic for URL creation 
* For every new URL, check if the oURL exist in key.
  * If the URL exist,
    * return the value from the database
  * If not
    * Generate the shortened value
    * Concate the generated value with the domain name and store into the database 

GET/{shortenedUrl}
The shortenedUrl will follow the regex /[a-zA-Z0-9]{9}/ (a string consisting of exactly 9 characters, where each character is either a lowercase letter (a-z), an uppercase letter (A-Z), or a digit (0-9).)
<<<<<<< HEAD

#### Sudo code logic for URL retrieving
* Search if the key exists in Cache
  * If the URL exist,
    * return the value from the database
  * If not
    * Search if the key exists in DynamoDB
      * If the URL exist,
        * return the value from the database
      * Else
        * return error with value not found
=======
>>>>>>> a507e6fe852165e2b30c4ec7a54c903f7fe55652

#### Sudo code logic for URL retrieving
* Search if the key exists in Cache
  * If the URL exist,
    * return the value from the database
  * If not
    * Search if the key exists in DynamoDB
      * If the URL exist,
        * return the value from the database
      * Else
        * return error with value not found
## Database schema design
We are levearaging the DynamoDB with 1 partition kay ,1 sort key. Since there is no complex relation for this application, there will only be one table "UrlMap". Here is the expected schema.

var schema = {
    TableName : "UrlMap",
    KeySchema : [ {
        AttributeName : "url",
        KeyType : "HASH"
    }, //Partition key
    {
        AttributeName : "shorternedUrl",
        KeyType : "RANGE"
    } //Sort key
    ],
    ProvisionedThroughput : {
        ReadCapacityUnits : 10,
        WriteCapacityUnits : 10
    }
};

Local Cli command
aws dynamodb create-table \
    --table-name UrlMap \
    --attribute-definitions \
        AttributeName=url,AttributeType=S \
        AttributeName=shortenedUrl,AttributeType=S \
    --key-schema \
        AttributeName=url,KeyType=HASH \
        AttributeName=shortenedUrl,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=10 \
    --table-class STANDARD \
    --endpoint-url http://localhost:8000

Command to list the table
aws dynamodb describe-table \
    --table-name UrlMap \
    --endpoint-url http://localhost:8000


Auto Scaling

Indexing


### The reason for this implementation

## Assumption
1. The same URL also returns with the same shortened URL
