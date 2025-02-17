# Use the Golang image from Docker Hub
FROM golang:1.24-bookworm AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Start a new stage from the golang image to reduce the size of the final image
FROM golang:1.24-bookworm

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY .env .env
# Expose port 8080 to be able to access the API
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
