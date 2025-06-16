# Build stage (cached)
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Install swag and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN swag init --generalInfo cmd/main.go --output docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fluently-backend ./cmd/main.go

# # Build stage (no-cache)
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
WORKDIR /app
COPY --from=builder /app/fluently-backend .
COPY --from=builder /app/docs ./docs
COPY .env .
EXPOSE ${APP_PORT}
CMD ["./fluently-backend"]
