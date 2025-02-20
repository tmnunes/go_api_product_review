
# GoLang API - Product Review

## Table of Contents
1. [The Challenge](#challenge)
2. [Architecture](#architecture)
3. [Running the App](#running)
4. [Testing the App](#testing)
5. [API Documentation](#API)
6. [Developing Preparation](#developing)

## The Challenge <a name="challenge"></a>
The goal is to create a RESTful API capable of managing products and product reviews. </br>
The API should calculate the average review score for each product.</br>

Technologies used: Go (Golang)

### Features:
- **Product API**: Supports CRUD operations and fetching by identifier.
    - Product information should return only the average product rating, not individual reviews.
- **Review API**: Supports Create, Update, Delete actions for reviews.
- **Review Notifications**: The service notifies the review service when a review is added, modified, or deleted.
- **Average Rating**: Every time a review is received, the average rating is recalculated and stored in persistent storage.
- **Caching**: Both product reviews and average ratings are cached.

## Architecture <a name="architecture"></a>
This project follows a Clean Architecture structure, with clear separation of concerns across models, services, and handlers to promote maintainability and scalability.

### Structure
```bash
/api               # Handles HTTP requests
/cache             # Caching mechanisms
/models            # Defines data structures
/service           # Contains business logic and data processing
/db                # Database connection and initialization
/cmd               # Entry point of the application and Swagger documentation
/middleware        # Authentication integration
/tests             # Unit tests
.env               # Environment variables configuration
Dockerfile         # Dockerfile for building the application image
docker-compose.yml # Docker Compose configuration for multi-container setup
```
### Use SQL DB or No-SQL
The options considered were PostgreSQL and MongoDB. </br> I chose PostgreSQL for its fixed schema and strong data integrity, which is essential for this project. </br>This choice ensures ACID compliance and data consistency across tables.

### Caching
The options was REDIS or Memcached. </br>
Redis was selected due to familiarity and its flexibility and scalability, making it a better fit for this project’s needs.

### Other possible aproaches
- **Microservices**: I opted not to use microservices to avoid unnecessary complexity in managing multiple smaller services. This approach simplifies the architecture and reduces overhead for the current scope.
- **Kubernetes**: Similarly, I avoided using Kubernetes due to its complexity and the resource costs associated with managing clusters. It would introduce unnecessary overhead for a project of this scale.

### Extras Added:
- **Authentication**: A simple Bearer token is used for API authentication, defined in the `.env` file.
- **Swagger Documentation**: Automatically generated API documentation for easy understanding of the API structure and interactions.
  ```bash
  swag init -g cmd/main.go
  ```

## Running the App <a name="running"></a>
You can run the service in two ways:

- **Docker (Recommended)**: This method builds and starts all dependencies and the solution. Ensure Docker is installed and running in the server/machine.
  ```bash
  docker-compose up --build
  ```

- **Standalone**: If you choose this option, make sure Redis and PostgreSQL are running and the configuration files are correctly set.
  ```bash
  go run cmd/main.go
  ```

## Testing the App <a name="testing"></a>
This project includes automated tests to ensure the correctness and reliability of the application.

### Tests Include:
- ✅ Unit testing for core application logic
- ✅ Mocking dependencies for isolated testing

### Running Tests
```bash
go test ./tests
```

## API Documentation <a name="API"></a>
This API provides endpoints for managing products and reviews. It uses Bearer token authentication for secure access.

You can view the full API documentation at: [Swagger UI](http://localhost:8080/swagger/index.html)

### API Endpoints

#### Products
- (GET) `/products`
- (POST) `/products`
- (GET) `/products/{id}`
- (PUT) `/products/{id}`
- (DELETE) `/products/{id}`

#### Reviews
- (POST) `/reviews`
- (PUT) `/reviews/{id}`
- (DELETE) `/reviews/{id}`

### Error Handling
HTTP Response Codes:
- **200 OK**: Successfully retrieved object(s)
- **201 Created**: Successfully created the object
- **204 No Content**: Successfully deleted the object
- **400 Bad Request**: Invalid input (e.g., missing required fields, invalid JSON)
- **401 Unauthorized**: User is not authenticated
- **403 Forbidden**: User does not have permission to perform the action
- **404 Not Found**: Object with the given ID does not exist
- **500 Internal Server Error**: Unexpected server error

## Developing Preparation <a name="developing"></a>

### Initial Setup for Local Machine
Run the following commands to prepare your local machine:
```bash
brew upgrade
brew install go
go version
docker version
```
For more details on installation, refer to:
- [Docker Installation Guide](https://docs.docker.com/desktop/setup/install/mac-install/)
- [DBeaver Installation](https://dbeaver.io/download/)

### Setting Up the Project and Docker
Install dependencies:
```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/gin-swagger/swaggerFiles
go get github.com/gorilla/mux
go get github.com/go-redis/redis/v8
go get -u gorm.io/driver/postgres
go get -u gorm.io/gorm
go get github.com/lib/pq
go get github.com/gin-contrib/static
```
