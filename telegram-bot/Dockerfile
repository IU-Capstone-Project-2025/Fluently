FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for Go module dependencies
RUN apk add --no-cache git

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies in a separate layer for better caching
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build the application with caching
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o fluently_bot ./cmd/main.go

FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and config from builder
COPY --from=builder /app/fluently_bot .
COPY --from=builder /app/config ./config

CMD ["./fluently_bot"]