## Installation
1. Clone the repository.
2. Run `go mod tidy` to install dependencies.
3. Start the application: `go run main.go`.

## API Endpoints
- `GET /start-stop`: Toggle automatic message sending (start/stop).
- `GET /sent-messages`: Retrieve a list of sent messages.

## Docker
Build and run the application using Docker:
```bash
docker-compose up --build
```
To view container logs:
```bash
docker logs message-sender-app-1
```
To stop and remove all containers:
```bash
docker-compose down
```

## MongoDB
Access MongoDB
To access the MongoDB instance running in Docker:
```bash
docker exec -it mongo mongosh
```

Add Data to MongoDB
Switch to the database:
```bash
use messageDB
```

Insert a document:
```bash
db.messages.insertOne({
content: "Hello, MongoDB!",
recipient: "+123456789",
sent: false
})
```

Verify the data:
```bash
db.messages.find()
```

## Troubleshooting
Common Issues
Application not accessible on port 8080:
```bash
- Ensure the application binds to 0.0.0.0:8080 in the code.
- Verify the docker-compose.yml file exposes port 8080.
```
MongoDB connection issues:
```bash
- Ensure the MongoDB URI is set to mongodb://mongo:27017 in the application.
```
Redis connection issues:
```bash
- Verify the Redis service is running and accessible on redis:6379.
```

## OpenAPI Specification
The OpenAPI specification for the API is available in the openapi.yaml file.
```bash
openapi: 3.0.0
info:
  title: Message Sender API
  description: API for managing message sending operations.
  version: 1.0.0
servers:
  - url: http://localhost:8080
    description: Local server
paths:
  /start-stop:
    get:
      summary: Toggle automatic message sending
      description: Starts or stops the automatic message sending process.
      responses:
        '200':
          description: Successfully toggled the state.
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: "Message sending started"
        '500':
          description: Internal server error
  /sent-messages:
    get:
      summary: Retrieve sent messages
      description: Fetches a list of all sent messages.
      responses:
        '200':
          description: A list of sent messages.
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: string
                      example: "64b7f3c2e4b0f5a1d2c3e4f5"
                    content:
                      type: string
                      example: "Hello, MongoDB!"
                    recipient:
                      type: string
                      example: "+123456789"
                    sent:
                      type: boolean
                      example: true
                    sentAt:
                      type: string
                      format: date-time
                      example: "2023-10-01T12:00:00Z"
        '404':
          description: No messages found
        '500':
          description: Internal server error
```