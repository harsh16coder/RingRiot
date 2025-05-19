# Use the official Golang image for development
FROM golang:1.24.3-alpine AS builder

# Set the working directory
WORKDIR /usr/src/gameserver

# Copy banned_words.txt into the image
COPY banned_words.txt /gameserver/banned_words.txt

# Copy dependency files and download modules
COPY go.mod go.sum ./ 
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Build the application binary
RUN go build -v -o /gameserver/main ./cmd/main.go

# Create a minimal image
FROM alpine:latest

# Create directory for the binary
RUN mkdir -p /gameserver

# Copy the binary from builder stage
COPY --from=builder /gameserver/main /gameserver/main

# Copy the .env file if needed (optional; usually handled via volumes or env_file)
# COPY .env .env

# Default command to run the application, referencing the .env file
CMD ["/gameserver/main", "--config", ".env"]
