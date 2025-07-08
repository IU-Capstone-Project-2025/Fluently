### Thesaurus API - Микросервис для рекомендаций по изучению лексики

Этот микросервис предоставляет персонализированные рекомендации по изучению английской лексики на основе известных пользователю слов, используя алгоритмы тематического подбора и частотного анализа.

> **Статус**: Интегрирован в проект Fluently с Docker поддержкой

## Основные возможности
- Рекомендации новых слов для изучения на основе известной лексики
- Учёт тематической связанности (topic/subtopic/subsubtopic)
- Фильтрация по уровню сложности (CEFR: A1-C2)
- Баланс между релевантностью и разнообразием тематик

## Требования
- Python 3.9+
- Установленные зависимости из `requirements.txt`
```bash
pip install -r requirements.txt
```

- Запустите сервис:
```bash
uvicorn thesaurus.app:app --reload --port 8000
```

Сервис будет доступен по адресу: `http://localhost:8000`

## Использование API

### Проверка работоспособности
```http
POST /health
Content-Type: application/json

{
  "ping": "test"
}
```

**Пример ответа:**
```json
{
  "status": "ok"
}
```

### Получение рекомендаций
```http
POST /api/recommend
Content-Type: application/json

{
  "words": ["apple", "bakery", "cook"]
}
```

**Параметры запроса:**
- `words`: Список известных слов (обязательный)
- `num_recommendations`: Количество рекомендаций (по умолчанию: 10)
- `max_cefr_level`: Максимальный уровень сложности (A1-C2, по умолчанию: C2)
- `max_per_subtopic`: Макс. слов из одной подтемы (по умолчанию: 2)

**Пример ответа:**
```json
[
  {
    "word": "pastry",
    "topic": "Food and drink",
    "subtopic": "Food",
    "subsubtopic": "Types of food",
    "CEFR_level": "B1",
    "score": 0.92
  },
  {
    "word": "cuisine",
    "topic": "Food and drink",
    "subtopic": "Cooking",
    "subsubtopic": "Cooking methods",
    "CEFR_level": "B2",
    "score": 0.89
  }
]
```

## Конфигурация

### Переменные окружения
Thesaurus API использует переменные окружения из корневого файла `.env` проекта:

```ini
# Настройки CORS для Thesaurus API (разделять запятыми)
THESAURUS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8070,https://fluently-app.ru

# Другие переменные из основного проекта...
```

### Docker окружение
API автоматически настраивается для работы в Docker-контейнере:
- Порт: 8002
- Автоматическая загрузка данных из `result.csv`
- CORS настраивается через переменную `THESAURUS_ALLOWED_ORIGINS`
- Интегрируется с основной сетью проекта `fluently_network`

## Структура проекта
```
├── thesaurus/
│   ├── app.py              # Основной код сервиса
│   └── requirements.txt    # Зависимости
├── data/
│   ├── oxford_vocabulary.csv  # Исходные данные
│   └── result.csv          # Обработанный датасет
├── .env.example            # Шаблон конфигурации
└── README.md               # Документация
```

## Разработка

### Локальная разработка
1. Активируйте виртуальное окружение:
```bash
python -m venv venv
source venv/bin/activate
```

2. Установите зависимости для разработки:
```bash
pip install -r requirements.txt
```

3. Запустите сервис в режиме разработки:
```bash
uvicorn app:app --reload --port 8002
```

### Docker разработка
Для разработки в составе всего проекта используйте:

```bash
# Локальная разработка со сборкой
docker compose -f docker-compose-local.yml up thesaurus-api

# Продакшн версия с готовым образом
docker compose up thesaurus-api
```

API будет доступен по адресу: `http://localhost:8002`

## Тестирование
Пример запроса с помощью cURL:
```bash
# Health check
curl -X POST http://localhost:8002/health \
  -H "Content-Type: application/json" \
  -d '{"ping": "test"}'

# Получение рекомендаций
curl -X POST http://localhost:8002/api/recommend \
  -H "Content-Type: application/json" \
  -d '{"words": ["computer", "software", "data"]}'
```

## Интеграция

Thesaurus API интегрируется с основным проектом Fluently:
- **Backend**: Может вызывать API для получения рекомендаций слов
- **Frontend**: Прямые запросы к API через настроенный CORS
- **Docker**: Автоматически развертывается как часть общего стека
- **CI/CD**: Автоматическая сборка и развертывание через GitHub Actions