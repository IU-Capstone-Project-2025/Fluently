FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fluently_bot ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/fluently_bot .
COPY --from=builder /app/config ./config

CMD ["./fluently_bot"]