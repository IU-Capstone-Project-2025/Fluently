[![Deploy](https://github.com/FluentlyOrg/Fluently-fork/actions/workflows/deploy.yml/badge.svg)](https://github.com/FluentlyOrg/Fluently-fork/actions/workflows/deploy.yml)

# Fluently - AI-Powered English Learning Platform
![Fluently Log](frontend-website/logo.jpg)
## Overview
Fluently is a cross-platform, open-source educational system for learning English, designed as a capstone project. It features:
- 🗣️ Realistic AI conversations
- 📚 Personalized vocabulary lessons
- 📈 Progress tracking
- 🧠 Adaptive learning based on user goals

**Platforms:** Android, iOS, Telegram Bot (soon), Web (Swagger UI for API)

## Links

Main project site (still in development)  
https://fluently-app.ru

Swagger (Fluently API documentation)  
https://fluently-app.ru/swagger/index.html

## Tech Stack
| Component       | Technologies                                                                 |
|-----------------|-------------------------------------------------------------------------------|
| Backend         | Go 1.24, Chi Router, GORM, PostgreSQL, Redis, Zap Logging, Swagger           |
| Mobile          | iOS (Swift), Android (Kotlin)                                                |
| Telegram Bot    | Go, Redis                                                                     |
| Infrastructure  | Docker, Docker Compose, Nginx, Let's Encrypt                                  |

---

## Installation & Testing

Fluently can be installed in two ways:

### 1. [Local/Development Installation](docs/Install_Local.md)
- **Recommended for teaching assistants and quick testing.**
- No domain or SSL required.
- All services run on `localhost` using Docker Compose.
- Test API, Swagger UI, and frontend separately.

### 2. [Full Production Installation](docs/Install_Full.md)
- **For advanced users or production deployment.**
- Requires your own domain and SSL certificates.
- Replicates the production environment.

---

## Documentation
- [Local Installation Guide](docs/Install_Local.md)
- [Full Production Installation Guide](docs/Install_Full.md)
- [Backend README](backend/README.md)

---

## License
[MIT](LICENSE)
