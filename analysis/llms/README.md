# Fluently LLM API

AI-powered conversation service for language learning using Groq and Gemini APIs.

## Features

- Support for multiple AI providers (Groq, Gemini)
- Automatic failover between providers
- API key rotation and cooldown management
- Prometheus metrics integration
- Health checks and status monitoring

## Environment Variables

The following environment variables should be set in the root `.env` file:

```bash
# Groq API Keys (comma-separated for multiple keys)
GROQ_API_KEYS=your_groq_key_1,your_groq_key_2

# Gemini API Keys (comma-separated for multiple keys)
GEMINI_API_KEYS=your_gemini_key_1,your_gemini_key_2
```

## API Endpoints

### Health Check
- `GET /health` - Check service health
- `GET /status` - Get detailed service status

### Chat Completion
- `POST /chat` - Full conversation endpoint
- `POST /chat/simple` - Simple single-message endpoint

### Monitoring
- `GET /metrics` - Prometheus metrics

## Example Usage

```bash
# Simple chat
curl -X POST "http://localhost:8003/chat/simple" \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, how are you?", "model_type": "fast"}'

# Full conversation
curl -X POST "http://localhost:8003/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {"role": "user", "content": "Hello, I want to learn English"}
    ],
    "model_type": "balanced"
  }'
```

## Development

To run locally:

```bash
pip install -r requirements.txt
python main.py
```

## Docker

Build and run with Docker:

```bash
docker build -t fluently-llm-api .
docker run -p 8003:8003 --env-file .env fluently-llm-api
```
