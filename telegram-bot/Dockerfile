# Production Dockerfile for Telegram Bot
FROM golang:1.24-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies in a separate layer for better caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build the application with optimizations for production
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o fluently_bot ./cmd/main.go

RUN mkdir -p /tmp/tts && chmod 777 /tmp/tts

# Production stage
FROM scratch

# Copy ca-certificates from builder for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Create non-root user for security
USER 1000:1000

# Copy binary from builder
COPY --from=builder /app/fluently_bot /fluently_bot

# Copy config directory if it exists
COPY --from=builder /app/config /config

# Set working directory
WORKDIR /

# Expose port if needed (adjust as per your application)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/fluently_bot", "--health-check"] || exit 1

# Run the application
CMD ["/fluently_bot"]