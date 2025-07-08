# Fluently Telegram Bot

This is the Telegram bot component of the Fluently language learning app.

## Features

- FSM (Finite State Machine) based conversation flows
- Redis-backed state persistence
- Support for multiple learning flows:
  - CEFR Level assessment
  - Vocabulary learning with spaced repetition
  - Exercise flows for practice
- Settings management:
  - Words per day
  - Notification preferences
  - CEFR level selection
- Structured error handling with recovery paths

## Architecture

The bot is built with a state-based architecture that allows for complex conversation flows:

- **FSM**: Manages user states and transitions
- **Router**: Routes updates to appropriate handlers
- **Handlers**: Process commands and messages
- **TempData**: Stores temporary conversation data
- **API Client**: Communicates with the backend API

### State Management

The bot uses Redis to persist user states and temporary data across interactions. The FSM implementation supports:

- Structured state transitions
- State validation
- Error recovery
- Temporary data storage
- Sub-state tracking for complex flows

## Key Components

### FSM

- `states.go`: Defines all possible states and valid transitions
- `memory.go`: Implements Redis-based state persistence and temporary data storage

### Handlers

- `service.go`: Provides handler service with state-based routing
- `handlers.go`: Implements command and message handlers

### Router

- `router.go`: Routes Telegram updates to appropriate handlers based on type and state

### Bot

- `bot.go`: Main bot implementation with FSM integration

## State Flows

### Onboarding Flow
1. Start → Welcome → Method Explanation → Spaced Repetition → Questionnaire
2. Questionnaire → Various Question States → CEFR Test

### CEFR Test Flow
1. Vocabulary Test → Test Groups (1-5) → Processing → Results → Level Determination

### Learning Flow
1. Lesson Start → Lesson In Progress → Show Words → Exercises → Lesson Complete

### Settings Flow
1. Settings → Various Setting States (Words Per Day, Notifications, CEFR Level)

## Running the Bot

```bash
# Development
go run cmd/main.go -debug -config=./config/config.dev.yaml

# Production
go run cmd/main.go -config=./config/config.prod.yaml
```

## Configuration

The bot requires the following configuration:

```yaml
bot:
  token: "your-telegram-bot-token"
  debug: false

api:
  base_url: "https://api.example.com/v1"

redis:
  address: "localhost:6379"
  db: 0
  password: ""
```

## 🚀 Features

- **Comprehensive Learning Flow**: 5-stage learning process from onboarding to lesson completion
- **Spaced Repetition System**: Scientific approach to vocabulary memorization  
- **Webhook-based**: High-performance webhook processing with async handling
- **Redis State Management**: Fast user progress and session storage
- **Background Tasks**: Asynq for scheduling lessons and notifications
- **API Integration**: Seamless connection with backend learning platform
- **Rate Limiting**: Protection against abuse and overload
- **Health Monitoring**: Built-in health checks and metrics

## 📋 Learning Flow

### Stage 1: Onboarding
- Welcome message and motivation
- Method explanation (spaced repetition)
- Scientific backing (Ebbinghaus forgetting curve)

### Stage 2: Personalization  
- User questionnaire (goals, confidence, habits)
- Vocabulary level assessment (5 groups of words)
- CEFR level determination
- Personalized learning plan creation

### Stage 3: Lesson Flow
- 10 words per lesson
- First block of 3 words with examples
- Individual word presentation with integrated repetition
- Audio and translation exercises

### Stage 4: Exercises
- Audio dictation (speech recognition)
- Translation checks (active recall)
- Immediate feedback and correction
- Progress tracking

### Stage 5: Progress Management
- Daily statistics and streaks
- Lesson completion rewards
- Scheduled reminders and notifications

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Telegram API  │───▶│  Webhook Server │───▶│ Handler Service │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │ Rate Limiter    │    │ FSM State Mgmt  │
                       │ Middleware      │    │ (Redis)         │
                       └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                              ┌─────────────────┐    ┌─────────────────┐
                              │ Background      │    │ Backend API     │
                              │ Tasks (Asynq)   │    │ Client          │
                              └─────────────────┘    └─────────────────┘
```

## 🛠️ Setup

### Prerequisites

- Go 1.21+
- Redis 6.0+
- Backend API running (Fluently backend)
- Telegram Bot Token

### Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd telegram-bot
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Set up Redis** 
```bash
# Using Docker
docker run -d --name redis -p 6379:6379 redis:7-alpine

# Or install locally
# Ubuntu/Debian
sudo apt-get install redis-server

# macOS
brew install redis
```

5. **Configure Telegram Bot**
   - Create bot with [@BotFather](https://t.me/BotFather)
   - Get bot token and add to `.env`
   - Set webhook URL (must be HTTPS in production)

### Configuration

Key environment variables in `.env`:

```bash
# Required
BOT_TOKEN=your_telegram_bot_token
WEBHOOK_URL=https://yourdomain.com/webhook
API_BASE_URL=https://your-backend-api.com

# Optional but recommended
WEBHOOK_SECRET=random_secret_string
REDIS_ADDR=localhost:6379
```

## 🚦 Running

### Development
```bash
go run cmd/main.go
```

### Production with Docker
```bash
docker-compose up -d
```

### Setting Webhook

The bot automatically configures webhooks, but you can also set manually:

```bash
curl -F "url=https://yourdomain.com/webhook" \
     -F "secret_token=your_webhook_secret" \
     https://api.telegram.org/bot<TOKEN>/setWebhook
```

## 📝 Project Structure

```
telegram-bot/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── api/
│   │   └── client.go          # Backend API client
│   ├── bot/
│   │   ├── fsm/
│   │   │   ├── states.go      # FSM state definitions
│   │   │   └── memory.go      # User progress tracking
│   │   └── handlers/
│   │       └── service.go     # Message handlers
│   ├── tasks/
│   │   ├── scheduler.go       # Asynq task scheduling
│   │   └── handlers.go        # Background task handlers
│   └── webhook/
│       └── server.go          # Webhook HTTP server
├── pkg/
│   └── logger/
│       └── logger.go          # Structured logging
├── docker-compose.yml         # Docker setup
├── Dockerfile                # Container image
└── .env.example              # Configuration template
```

## 🔄 FSM States

The bot uses a comprehensive FSM to manage user learning progress:

### Onboarding States
- `StateStart` - Initial state
- `StateWelcome` - Welcome message
- `StateMethodExplanation` - Learning method explanation
- `StateSpacedRepetition` - Spaced repetition explanation

### Assessment States  
- `StateQuestionnaire` - User questionnaire
- `StateVocabularyTest` - Vocabulary level test
- `StateLevelDetermination` - CEFR level determination

### Learning States
- `StateLessonStart` - Lesson beginning
- `StateShowingWords` - Word presentation
- `StateExerciseReview` - Exercise phase
- `StateLessonComplete` - Lesson completion

### Exercise States
- `StateAudioDictation` - Audio exercises
- `StateTranslationCheck` - Translation exercises
- `StateWaitingForAudio` - Awaiting audio response
- `StateWaitingForTranslation` - Awaiting translation

## 🔧 API Integration

The bot integrates with the Fluently backend API for:

- **User Authentication**: Telegram ↔ Google account linking
- **Lesson Generation**: Dynamic lesson creation based on user level
- **Progress Tracking**: Real-time progress synchronization
- **Content Delivery**: Words, sentences, audio, exercises

### Account Linking

Users can link Telegram accounts with Google accounts:

1. Bot calls `/api/v1/telegram/create-link`
2. User clicks magic link → Google OAuth
3. Backend links accounts
4. Bot gets confirmation via `/api/v1/telegram/check-status`

## ⚡ Performance Features

### Webhook Processing
- **Async Processing**: Immediate 200 OK response to Telegram
- **Goroutine Workers**: Parallel message processing
- **Rate Limiting**: 100 req/sec with 200 burst
- **Request Size Limits**: 1MB max body size

### Redis Optimization
- **Connection Pooling**: Efficient Redis connections
- **TTL Management**: Automatic data expiration
- **Compressed Storage**: JSON compression for large data
- **Session Cleanup**: Hourly cleanup of expired data

### Background Tasks
- **Lesson Reminders**: Scheduled daily reminders
- **Progress Sync**: Periodic API synchronization  
- **Notifications**: Daily facts and motivation
- **Cleanup Tasks**: Data maintenance and optimization

## 📊 Monitoring

### Health Checks
- `GET /health` - Basic health status
- `GET /ready` - Service readiness (Redis, API connectivity)
- `GET /metrics` - Basic metrics (requests, errors, uptime)

### Logging
- **Structured Logging**: JSON format with Zap
- **Log Levels**: Debug, Info, Warn, Error
- **Context Tracking**: Request IDs and user tracking
- **Error Reporting**: Detailed error context

## 🐳 Docker Deployment

### Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f bot

# Stop services
docker-compose down
```

### Production Deployment
1. Set up reverse proxy (nginx/Caddy) for HTTPS
2. Configure SSL certificates
3. Set production environment variables
4. Enable log rotation
5. Set up monitoring (Prometheus/Grafana)

## 🔐 Security

### Webhook Security
- **Secret Token Validation**: Telegram webhook secret
- **HMAC Signature Verification**: Request signature validation
- **Rate Limiting**: Protection against flooding
- **Request Size Limits**: Prevent large payload attacks

### Data Protection
- **Redis Authentication**: Password-protected Redis
- **API Key Management**: Secure API key storage
- **Session Expiration**: Automatic cleanup of user sessions
- **Input Validation**: Sanitization of user inputs

## 🐛 Troubleshooting

### Common Issues

**Bot not responding to messages:**
```bash
# Check webhook status
curl https://api.telegram.org/bot<TOKEN>/getWebhookInfo

# Check logs
docker-compose logs bot
```

**Redis connection issues:**
```bash
# Test Redis connectivity
redis-cli ping

# Check Redis logs
docker-compose logs redis
```

**API integration problems:**
```bash
# Test API connectivity
curl -H "X-API-Key: your_key" https://your-api.com/health

# Check API client logs
grep "API error" logs/bot.log
```

### Debug Mode

Enable debug logging:
```bash
LOG_LEVEL=debug go run cmd/main.go
```

## 📈 Scaling

### Horizontal Scaling
- **Multiple Bot Instances**: Load balancing with shared Redis
- **Worker Separation**: Dedicated Asynq workers
- **Database Sharding**: User-based Redis sharding

### Performance Optimization
- **Connection Pooling**: Redis and HTTP connection pools
- **Caching**: Frequent data caching in Redis
- **Batch Processing**: Bulk API operations
- **Async Operations**: Non-blocking task processing

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Telebot v3](https://gopkg.in/telebot.v3) - Telegram Bot framework
- [Asynq](https://github.com/hibiken/asynq) - Background task processing
- [Chi](https://github.com/go-chi/chi) - HTTP router
- [Redis](https://redis.io/) - In-memory data structure store
- [Zap](https://github.com/uber-go/zap) - Structured logging
