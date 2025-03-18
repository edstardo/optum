# Use official Golang image to build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.* ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -o service ./cmd/trader-service

# Use a minimal base image for the final stage
FROM ubuntu:22.04

# Set working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/service .

EXPOSE 8080

# Command to run the application
CMD ["./service"]