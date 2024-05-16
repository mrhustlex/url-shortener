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
Logic:
Assuming the NoSQL Database is a hashmap with original url(oUrl) as key and shortened URL(sURL) as value

#### Sudo code logic for URL creation 
* For every new URL, check if the oURL exist in key.
  * If the URL exist,
    * return the value from the database
  * If not
    * Generate the shortened value
    * Concate the generated value with the domain name and store into the database 

GET/{shortenedUrl}
* The shortenedUrl will follow the regex /[a-zA-Z0-9]{9}/ (a string consisting of exactly 9 characters, where each character is either a lowercase letter (a-z), an uppercase letter (A-Z), or a digit (0-9).)

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

### The reason for this implementation

## Assumption
1. The same URL also returns with the same shortened URL
