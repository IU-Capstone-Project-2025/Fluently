# FastAPI Distractor Generator Service

High-performance API for generating word distractors using BERT model, optimized for low latency with async processing and model caching.

## Features

- **Low Latency**: Model loaded once at startup, async processing
- **Clean Architecture**: Modular design with separation of concerns
- **Error Handling**: Comprehensive error handling and validation
- **Health Checks**: Built-in health monitoring
- **Docker Support**: Containerized deployment ready
- **API Documentation**: Auto-generated OpenAPI/Swagger docs

## Architecture

```
distractor_api/
├── main.py              # FastAPI app entry point
├── models/
│   └── schemas.py       # Pydantic models
├── services/
│   └── distractor_service.py  # Business logic
├── api/
│   └── routes.py        # API endpoints
├── requirements.txt     # Dependencies
└── Dockerfile          # Container setup
```

## API Endpoints

### POST /api/v1/generate-distractors

Generate distractor words for a target word in a sentence.

**Request:**
```json
{
  "sentence": "The cat caught the mouse in the kitchen",
  "word": "mouse"
}
```

**Response:**
```json
{
  "pick_options": ["mouse", "rat", "bird", "fish"]
}
```

### GET /health

Health check endpoint to verify service status.

**Response:**
```json
{
  "status": "healthy",
  "model_loaded": true,
  "timestamp": 1640995200.0
}
```

## Quick Start

### Local Development

1. **Install dependencies:**
   ```bash
   pip install -r requirements.txt
   ```

2. **Run the service:**
   ```bash
   cd backend
   python -m uvicorn distractor_api.main:app --host 0.0.0.0 --port 8001 --reload
   ```

3. **Access the API:**
   - API: http://localhost:8001
   - Docs: http://localhost:8001/docs
   - Health: http://localhost:8001/health

### Docker Deployment

1. **Build the image:**
   ```bash
   cd backend/distractor_api
   docker build -t distractor-api .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8001:8001 distractor-api
   ```

## Performance Optimizations

- **Model Caching**: BERT model loaded once at startup
- **Async Processing**: Non-blocking I/O operations
- **GPU Support**: Automatic CUDA detection and usage
- **Model Warmup**: Dummy inference during initialization
- **Single Worker**: Avoids model loading overhead
- **Optimized Inference**: Efficient tokenization and prediction

## API Usage Examples

### Python (requests)

```python
import requests

response = requests.post(
    "http://localhost:8001/api/v1/generate-distractors",
    json={
        "sentence": "The programmer fixed all bugs in the code",
        "word": "bugs"
    }
)

result = response.json()
print(result["pick_options"])  # ["bugs", "errors", "issues", "problems"]
```

### cURL

```bash
curl -X POST "http://localhost:8001/api/v1/generate-distractors" \
     -H "Content-Type: application/json" \
     -d '{
       "sentence": "Birds fly south for the winter",
       "word": "fly"
     }'
```

### JavaScript (fetch)

```javascript
const response = await fetch('http://localhost:8001/api/v1/generate-distractors', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    sentence: "She opened the window to breathe fresh air",
    word: "window"
  })
});

const result = await response.json();
console.log(result.pick_options);
```

## Error Handling

The API provides detailed error responses:

- **400 Bad Request**: Invalid input (empty sentence/word, word not in sentence)
- **422 Validation Error**: Pydantic validation failures
- **503 Service Unavailable**: Model not loaded or internal errors

## Configuration

### Environment Variables

- `PYTHONPATH`: Set to `/app` for proper module imports
- `PYTHONUNBUFFERED`: Set to `1` for immediate output
- Custom model configurations can be added to the service

### Hardware Requirements

- **CPU**: 2+ cores recommended
- **RAM**: 4GB+ (8GB+ for better performance)
- **GPU**: Optional, CUDA-compatible for faster inference
- **Storage**: 2GB+ for model caching

## Monitoring

- Health check endpoint: `/health`
- Request timing logged automatically
- Built-in Docker health checks
- Prometheus metrics can be added

## Development

### Code Style

- Type hints throughout
- Async/await patterns
- Comprehensive error handling
- Clean architecture principles

### Testing

Run tests with:
```bash
pytest tests/
```

### Extending

To add new features:
1. Update schemas in `models/schemas.py`
2. Implement logic in `services/`
3. Add routes in `api/routes.py`
4. Update documentation

## License

See LICENSE file for details. 