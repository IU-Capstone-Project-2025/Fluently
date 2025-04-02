<<<<<<< HEAD
<<<<<<< HEAD
# Fluently
=======
# go-backend
=======
## Setup and Usage
>>>>>>> d6ccba0 (Add initial project configuration files and dependencies)

### Requirements

- Go 1.23+
- PostgreSQL
- Redis

### 0. Git

Сделал ветку develop
От неё уже есть две ветки:
- feature/models, там пишешь код моделек
- feature/handlers - код хендлеров
- Можешь создавать по такому же принципу ветки и делать в них

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

<<<<<<< HEAD
Every project is different, so consider which of these sections apply to yours. The sections used in the template are suggestions for most open source projects. Also keep in mind that while a README can be too long and detailed, too long is better than too short. If you think your README is too long, consider utilizing another form of documentation rather than cutting out information.

## Name
Choose a self-explaining name for your project.

## Description
Let people know what your project can do specifically. Provide context and add a link to any reference visitors might be unfamiliar with. A list of Features or a Background subsection can also be added here. If there are alternatives to your project, this is a good place to list differentiating factors.

## Badges
On some READMEs, you may see small images that convey metadata, such as whether or not all the tests are passing for the project. You can use Shields to add some to your README. Many services also have instructions for adding a badge.

## Visuals
Depending on what you are making, it can be a good idea to include screenshots or even a video (you'll frequently see GIFs rather than actual videos). Tools like ttygif can help, but check out Asciinema for a more sophisticated method.

## Installation
Within a particular ecosystem, there may be a common way of installing things, such as using Yarn, NuGet, or Homebrew. However, consider the possibility that whoever is reading your README is a novice and would like more guidance. Listing specific steps helps remove ambiguity and gets people to using your project as quickly as possible. If it only runs in a specific context like a particular programming language version or operating system or has dependencies that have to be installed manually, also add a Requirements subsection.

## Usage
Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

## Support
Tell people where they can go to for help. It can be any combination of an issue tracker, a chat room, an email address, etc.

## Roadmap
If you have ideas for releases in the future, it is a good idea to list them in the README.

## Contributing
State if you are open to contributions and what your requirements are for accepting them.

For people who want to make changes to your project, it's helpful to have some documentation on how to get started. Perhaps there is a script that they should run or some environment variables that they need to set. Make these steps explicit. These instructions could also be useful to your future self.

You can also document commands to lint the code or run tests. These steps help to ensure high code quality and reduce the likelihood that the changes inadvertently break something. Having instructions for running tests is especially helpful if it requires external setup, such as starting a Selenium server for testing in a browser.

## Authors and acknowledgment
Show your appreciation to those who have contributed to the project.

## License
For open source projects, say how it is licensed.

## Project status
If you have run out of energy or time for your project, put a note at the top of the README saying that development has slowed down or stopped completely. Someone may choose to fork your project or volunteer to step in as a maintainer or owner, allowing your project to keep going. You can also make an explicit request for maintainers.
>>>>>>> 27352b1 (Initial commit)
=======
- [Chi Router](https://github.com/go-chi/chi): Lightweight, idiomatic HTTP router
- [GORM](https://gorm.io/): ORM library for Golang
- [Viper](https://github.com/spf13/viper): Configuration solution
- [Zap](https://github.com/uber-go/zap): Structured logging
- [Swaggo](https://github.com/swaggo/swag): Swagger 2.0 generator for Go
- [Air](https://github.com/cosmtrek/air): Live reload for Go apps
>>>>>>> d6ccba0 (Add initial project configuration files and dependencies)
