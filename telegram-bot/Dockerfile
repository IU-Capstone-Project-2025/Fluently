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

# Create /tmp/tts and set permissions in builder stage
RUN mkdir -p /tmp/tts && chmod 777 /tmp/tts

# Production stage
FROM alpine:latest

# Copy ca-certificates from builder for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy /tmp/tts directory from builder (optional, but we'll recreate it anyway)
# COPY --from=builder /tmp/tts /tmp/tts

# Create non-root user for security (ensure it exists in alpine)
RUN adduser -D -u 1000 appuser

# Ensure /tmp/tts exists and is writable by the app user
USER root
RUN mkdir -p /tmp/tts && chmod 777 /tmp/tts && chown appuser:appuser /tmp/tts

# Copy binary from builder
COPY --from=builder /app/fluently_bot /fluently_bot

# Copy config directory if it exists
COPY --from=builder /app/config /config

# Set working directory
WORKDIR /

# Expose port if needed (adjust as per your application)
EXPOSE 8080

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/fluently_bot", "--health-check"] || exit 1

# Run the application
CMD ["/fluently_bot"]