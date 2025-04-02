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

```main.go
	logger.Log.Info("Logger initialization successful!")
	logger.Log.Info("App starting",
		zap.String("name", config.GetAppName()),
		zap.String("address", config.GetAppHost()+":"+config.GetAppPort()),
		zap.String("dsn", config.GetPostgresDSN()),
	)
```

# Project Structure
## 🗂️ Project Structure — `fluently/go-backend`

```txt
.
├── cmd/                            # Точка входа в приложение
│   └── main.go                     # Запуск HTTP-сервера, зависимостей и маршрутов
├── docs/                           # Swagger-документация (сгенерировано через swag)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod, go.sum                  # Зависимости проекта (модуль Go)
├── internal/                       # Основная бизнес-логика (handlers, сервисы, доступ к данным)
│   ├── api/
│   │   └── v1/
│   │       ├── handlers/           # HTTP-обработчики (controllers)
│   │       │   └── *.go            # Например: word_handler.go, user_handler.go и т.д.
│   │       └── routes/             # Роутинг chi (RegisterWordRoutes, и т.п.)
│   │           └── *.go
│   ├── config/                     # Загрузка конфигурации (viper)
│   │   └── config.go
│   ├── db/                         # Инициализация базы, миграции, подключения (ещё пусто)
│   ├── repository/                 # Слой доступа к данным (models, postgres-реализации, DTO)
│   │   ├── models/                 # GORM-модели таблиц
│   │   ├── postgres/               # Реализации репозиториев через GORM
│   │   └── schemas/                # DTO-схемы (вход/выход)
│   ├── router/                     # Сборка chi.Router
│   │   └── router.go
│   ├── swagger/                    # Связь между swagger-доками и chi (опционально)
│   └── utils/                      # Хелперы, утилиты, форматирование, ошибки и т.д.
├── migrations/                     # SQL- или go-модули для миграций базы данных
├── pkg/
│   └── logger/                     # Zap-логгер (переиспользуемый)
│       └── logger.go
└── README.md                       # Главный файл описания проекта
```

---

## Общая концепция

- `internal/` — основная логика проекта, разбитая по слоям
- `repository/` — реализация работы с БД: модели, схемы и репозитории
- `api/v1/` — REST API (обработчики + маршруты)
- `pkg/` — внешний код, пригодный для повторного использования

## Dependencies

- [Chi Router](https://github.com/go-chi/chi): Lightweight, idiomatic HTTP router
- [GORM](https://gorm.io/): ORM library for Golang
- [Viper](https://github.com/spf13/viper): Configuration solution
- [Zap](https://github.com/uber-go/zap): Structured logging
- [Swaggo](https://github.com/swaggo/swag): Swagger 2.0 generator for Go
- [Air](https://github.com/cosmtrek/air): Live reload for Go apps