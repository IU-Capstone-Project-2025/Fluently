## Setup and Usage

### Requirements

- Go 1.20+
- PostgreSQL
- Redis

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Run in Development Mode

```bash
air
```

> Убедись, что у тебя установлен `air`. Если нет — установи:
> ```bash
> go install github.com/air-verse/air@latest
> ```

### 3. Generate Swagger Docs

```bash
swag init --generalInfo cmd/main.go --output docs
```

> Убедись, что у тебя установлен `swag`:
> ```bash
> go install github.com/swaggo/swag/cmd/swag@latest
> ```

Swagger-документация будет доступна по маршруту `/swagger/index.html`, если подключён `httpSwagger.Handler`.

### 4. Example of logging


---

## Dependencies

- [Chi Router](https://github.com/go-chi/chi): Lightweight, idiomatic HTTP router
- [GORM](https://gorm.io/): ORM library for Golang
- [Viper](https://github.com/spf13/viper): Configuration solution
- [Zap](https://github.com/uber-go/zap): Structured logging
- [Swaggo](https://github.com/swaggo/swag): Swagger 2.0 generator for Go
- [Air](https://github.com/cosmtrek/air): Live reload for Go apps