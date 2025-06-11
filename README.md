
# Fluently - AI-Powered English Learning Platform

## Overview
Fluently addresses the gap in effective English learning tools for A2-C2 CEFR level learners by offering:
- 🗣️ Realistic AI conversations
- 📚 Personalized vocabulary lessons
- 📈 Progress tracking across devices
- 🧠 Adaptive learning based on user goals and schedule

## Tech Stack
### Core Components
| Component       | Technologies                                                                 |
|-----------------|-------------------------------------------------------------------------------|
| Backend         | Go 1.24, Chi Router, GORM, PostgreSQL, Redis, Zap Logging, Swagger           |
| Mobile          | iOS (Swift), Android (Kotlin)                                                |
| Telegram Bot    | Go, Redis                                                                     |
| Infrastructure  | Docker, Docker Compose, Nginx, Let's Encrypt                                  |

## Getting Started
### Prerequisites
- Docker 20.10+
- Docker Compose 2.20+
- Git 2.40+
### Deployment Steps
 - Clone repository:
```bash
git clone https://github.com/IU-Capstone-Project-2025/Fluently.git
cd fluently
```
- Create .env file in backend directory:
```conf
# backend/.env
BOT_TOKEN=<your-telegram-bot-token>
APP_NAME=Fluently_prod
APP_HOST=0.0.0.0
APP_PORT=8080

DB_USER=fluently_user
DB_PASSWORD=secure_password
DB_HOST=postgres
DB_PORT=5432
DB_NAME=fluently_prod
```

- Update domain in Nginx config:
```bash
# Edit backend/swagger/nginx.conf
server_name your-domain.com; # Replace swagger.fluently-app.ru
```
 - Start services:
```bash
docker compose -f backend/docker-compose.yml up -d --build
```
 - CI/CD Setup (Optional)
For automated deployments, add these GitHub Secrets:
```yaml
DEPLOY_HOST: SSH server IP
DEPLOY_USERNAME: SSH username
DEPLOY_SSH_KEY: Private SSH key
```
- Example workflow:
```yaml
# .github/workflows/deploy.yml
name: Deploy
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.DEPLOY_HOST }}
          username: ${{ secrets.DEPLOY_USERNAME }}
          key: ${{ secrets.DEPLOY_SSH_KEY }}
          script: |
            cd /home/deploy/Fluently
            git pull && docker compose up -d --build
```
