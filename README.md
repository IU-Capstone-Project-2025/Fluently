[![Deploy](https://github.com/FluentlyOrg/Fluently-fork/actions/workflows/deploy.yml/badge.svg)](https://github.com/FluentlyOrg/Fluently-fork/actions/workflows/deploy.yml)

# Fluently - AI-Powered English Learning Platform
![Fluently Log](frontend-website/logo-t.png)
## Overview
Fluently is a cross-platform, open-source educational system for learning English, designed as a capstone project. It features:
- 🦾 Interactive Chat: Practice conversations with AI
- 📚 Personalized vocabulary lessons
- 📈 Progress tracking
- 🧠 Adaptive learning based on user goals

**Platforms:** Android, iOS, Telegram Bot

## Links

Main project site 
https://fluently-app.ru

Telegram Bot 
http://t.me/FluentlyInEnglishBot

Terms of Use
https://fluently-app.ru/terms

### Latest Releases

[Android](https://github.com/FluentlyOrg/Fluently-fork/releases/download/v1.0.0-mob/app-release.apk)

[iOS](https://github.com/FluentlyOrg/Fluently-fork/releases/download/v1.0.0-mob/Fluently.ipa)

## Tech Stack
| Component       | Technologies                                                                 |
|-----------------|-------------------------------------------------------------------------------|
| Backend         | Go 1.24, Chi Router, GORM, PostgreSQL, Redis, Zap Logging, Swagger           |
| Mobile          | iOS (Swift), Android (Kotlin)                                                |
| Telegram Bot    | Go, Redis                                                                     |
| Infrastructure  | Docker, Docker Compose, Nginx, Prometheus, Grafana, Loki, PostgreSQL, Redis, Cloudflare |

---

## Installation & Testing

> [!IMPORTANT]
> Full and local installations are **only** supported on Linux (Ubuntu 22.04+).

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
- [Platform Support & Troubleshooting](docs/Platform_Support.md)
- [Backend README](backend/README.md)

---

## License
[MIT](LICENSE)
