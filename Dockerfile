# Build stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main cmd/bot/main.go

# Final stage - use debian slim for better compatibility
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /data

# Copy the binary from builder
COPY --from=builder /app/main /usr/local/bin/spends-bot

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep -f spends-bot || exit 1

CMD ["spends-bot"]
