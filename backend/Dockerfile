# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build with separate cache mounts
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o /app/fluently-backend ./cmd/main.go


# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fluently-backend .
COPY --from=builder /app/docs ./docs
COPY .env .
EXPOSE ${APP_PORT}
CMD ["./fluently-backend"]