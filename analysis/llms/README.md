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

## Troubleshooting

### Common Issues

1. **"TypeError: Client.__init__() got an unexpected keyword argument 'proxies'"**
   - **Solution**: Update the `groq` library to version 0.11.0 or later
   - **Cause**: Older versions of the Groq library have compatibility issues

2. **"All providers failed"**
   - **Solution**: Check that API keys are properly set in environment variables
   - **Solution**: Verify API keys are valid and have sufficient quota
   - **Solution**: Check network connectivity to API endpoints

3. **Pydantic warnings about "model_" namespace**
   - **Solution**: This is handled by setting `model_config = {"protected_namespaces": ()}` in the models
   - **Note**: These warnings don't affect functionality

4. **"Application startup failed"**
   - **Solution**: Check that all required environment variables are set
   - **Solution**: Verify API keys are valid
   - **Note**: The service will still start even if AI initialization fails

### Logs

The service logs important events including:
- API key initialization
- Request failures  
- Rate limiting events
- Provider fallbacks

Set `LOG_LEVEL=DEBUG` for detailed logging.

### Testing

Run the test script to verify functionality:
```bash
python test_api.py
```
