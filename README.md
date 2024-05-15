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
Database: Redis, NoSQL(Amazon DynamoDB)

## API design
POST/new_url

GET/{shortenedUrl}
* The shortenedUrl will follow the regex /[a-zA-Z0-9]{9}/ (a string consisting of exactly 9 characters, where each character is either a lowercase letter (a-z), an uppercase letter (A-Z), or a digit (0-9).)

## Database schema design

### The reason for this implementation

## Assumption
1. The same URL also returns with the same shortened URL
2. 
