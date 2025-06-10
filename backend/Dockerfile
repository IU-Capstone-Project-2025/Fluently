# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN \
--mount=type=cache,target=/go/pkg/mod \
CGO_ENABLED=0 GOOS=linux go build -o /app/fluently-backend ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fluently-backend .
COPY --from=builder /app/docs ./docs
COPY .env .
EXPOSE ${APP_PORT}
CMD ["./fluently-backend"]