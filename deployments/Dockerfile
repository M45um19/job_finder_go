# Stage 1: Build the Go application
FROM golang:1.25.6-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Stage 2: Run the application
FROM alpine:latest

# Install CA certificates to allow secure SSL connections to the database
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 5000

# Run the app
CMD ["./main"]
