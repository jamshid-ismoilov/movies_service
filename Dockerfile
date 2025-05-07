# Multi-stage build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build tools and dependencies
RUN apk add --no-cache git

# Copy go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install sql-migrate CLI tool
RUN go install github.com/rubenv/sql-migrate/...@latest

# Copy the source code
COPY . .

# Build the application
RUN go build -o movies_service .

# Final stage for running
FROM alpine:3.17

WORKDIR /app

# Copy the compiled binary and required files
COPY --from=builder /app/movies_service /app/
COPY --from=builder /go/bin/sql-migrate /usr/bin/
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/dbconfig.yml /app/

# Copy entrypoint script
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
