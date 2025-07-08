# Build stage (optimized with caching)

FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install git and envsubst (required for some Go modules and env substitution)
RUN apk add --no-cache git gettext

# Copy dependency files first for better caching
COPY go.mod go.sum ./

# Download dependencies with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Install swag with cache mount
RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code (this layer will be rebuilt when code changes)
COPY . .

# Substitute environment variables in main.go and generate docs
ARG SWAGGER_HOST
RUN envsubst '${SWAGGER_HOST}' < cmd/main.go > cmd/main.go.tmp && mv cmd/main.go.tmp cmd/main.go
RUN swag init --generalInfo cmd/main.go --output docs

# Build the application with build cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /app/fluently-backend ./cmd/main.go


# Build stage (no-cache)

# FROM golang:1.24-alpine AS builder
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# # Install swag and generate docs
# RUN go install github.com/swaggo/swag/cmd/swag@latest
# COPY . .
# RUN swag init --generalInfo cmd/main.go --output docs
# RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fluently-backend ./cmd/main.go

# Final stage
FROM alpine:latest

# Install curl for health checks
RUN apk add --no-cache curl

WORKDIR /app
COPY --from=builder /app/fluently-backend .
COPY --from=builder /app/docs ./docs

# Create a simple startup script that handles environment setup
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "🚀 Starting Fluently Backend..."' >> /app/start.sh && \
    echo 'echo "📍 Working directory: $(pwd)"' >> /app/start.sh && \
    echo 'echo "📋 Environment variables loaded"' >> /app/start.sh && \
    echo 'exec ./fluently-backend "$@"' >> /app/start.sh && \
    chmod +x /app/start.sh

# Set default port
EXPOSE 8070

CMD ["/app/start.sh"]
