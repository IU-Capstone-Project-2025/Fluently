### Thesaurus API - Микросервис для рекомендаций по изучению лексики

Этот микросервис предоставляет персонализированные рекомендации по изучению английской лексики на основе известных пользователю слов, используя алгоритмы тематического подбора и частотного анализа.

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
uvicorn backend.app:app --reload --port 8000
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
Создайте файл `.env` в корне проекта:
```ini
# Настройки CORS (разделять запятыми)
ALLOWED_ORIGINS=http://localhost:3000,https://your-frontend.com

# Параметры рекомендаций по умолчанию
DEFAULT_RECOMMENDATIONS=15
DEFAULT_MAX_SUBTOPIC=3
```

### Настройка CORS
Измените разрешённые домены в `backend/app.py`:
```python
app.add_middleware(
    CORSMiddleware,
    allow_origins=os.getenv("ALLOWED_ORIGINS", "http://localhost:3000").split(","),
    ...
)
```

## Структура проекта
```
├── backend/
│   ├── app.py              # Основной код сервиса
│   └── requirements.txt    # Зависимости
├── data/
│   ├── oxford_vocabulary.csv  # Исходные данные
│   └── result.csv          # Обработанный датасет
├── .env.example            # Шаблон конфигурации
└── README.md               # Документация
```

## Разработка

1. Активируйте виртуальное окружение:
```bash
python -m venv venv
source venv/bin/activate
```

2. Установите зависимости для разработки:
```bash
pip install -r backend/requirements.txt
```

3. Запустите сервис в режиме разработки:
```bash
uvicorn backend.app:app --reload --port 8000
```

## Тестирование
Пример запроса с помощью cURL:
```bash
curl -X POST http://localhost:8000/api/recommend \
  -H "Content-Type: application/json" \
  -d '{"words": ["computer", "software", "data"]}'
```