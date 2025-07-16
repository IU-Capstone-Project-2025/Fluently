# Full Production Installation Guide

This guide explains how to deploy Fluently in a production environment with your own domain, SSL certificates, and all services running as in production.

---

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/) (v2.0+)
- [Git](https://git-scm.com/)
- Your own domain name (e.g., `your-domain.com`)
- VPS or cloud server with public IP and at least 8GB RAM
- Basic understanding of DNS configuration
- SSH access to your server

---

## 1. Clone the Repository
```bash
git clone https://github.com/FluentlyOrg/Fluently-fork.git
cd Fluently-fork
```

---

## 2. Create Environment File
Create a `.env` file in the project root with the following configuration:

```env
# ===========================================
# APPLICATION CONFIGURATION
# ===========================================

# App Configuration
APP_NAME=FluentlyAPI
APP_HOST=0.0.0.0
APP_PORT=8070

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_database_password_here
DB_NAME=postgres

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here_make_it_long_and_secure_minimum_32_characters
JWT_EXPIRATION=24h
REFRESH_EXPIRATION=720h

# Google OAuth Configuration
IOS_GOOGLE_CLIENT_ID=your_ios_google_client_id
ANDROID_GOOGLE_CLIENT_ID=your_android_google_client_id
WEB_GOOGLE_CLIENT_ID=your_web_google_client_id
WEB_GOOGLE_CLIENT_SECRET=your_web_google_client_secret

# Logging Configuration
LOG_LEVEL=info
LOG_PATH=./logs/app.log

# Security Configuration
PASSWORD_MIN_LENGTH=8
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1h

# Swagger Configuration
SWAGGER_ALLOWED_EMAILS=admin@yourdomain.com

# ===========================================
# INFRASTRUCTURE CONFIGURATION
# ===========================================

# Domain Configuration
DOMAIN=yourdomain.com
CERT_NAME=yourdomain

# ZeroTier Configuration (for secure admin access)
ZEROTIER_IP=your.zerotier.ip

# ===========================================
# ADMIN SERVICES CONFIGURATION
# ===========================================

# Directus CMS Configuration
DIRECTUS_ADMIN_EMAIL=admin@yourdomain.com
DIRECTUS_ADMIN_PASSWORD=your_super_secure_directus_password_here
DIRECTUS_SECRET_KEY=your_directus_secret_key_here_make_it_long_and_secure

# Grafana Configuration
GRAFANA_ADMIN_PASSWORD=your_super_secure_grafana_password_here

# ===========================================
# TELEGRAM BOT CONFIGURATION
# ===========================================

# Telegram Bot Configuration
BOT_TOKEN=your_telegram_bot_token_here
WEBHOOK_URL=https://yourdomain.com/webhook
WEBHOOK_PATH=/webhook
WEBHOOK_SECRET=your_webhook_secret_here

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Backend API Configuration
API_BASE_URL=https://yourdomain.com/api/v1
API_TIMEOUT=30

# Asynq Configuration (for background tasks)
ASYNQ_REDIS_ADDR=localhost:6379
ASYNQ_REDIS_PASSWORD=
ASYNQ_REDIS_DB=1
ASYNQ_CONCURRENCY=10

# Environment
ENVIRONMENT=production

# ===========================================
# LLM API CONFIGURATION
# ===========================================

# AI/LLM API Keys (comma-separated for multiple keys)
# Get Groq API keys from: https://console.groq.com/keys
# Get Gemini API keys from: https://makersuite.google.com/app/apikey
GROQ_API_KEYS=your_groq_api_key_here,your_second_groq_api_key_here
GEMINI_API_KEYS=your_gemini_api_key_here,your_second_gemini_api_key_here

# ===========================================
# EXTERNAL SERVICES
# ===========================================

# Thesaurus API Configuration
THESAURUS_ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```
---

## 3. Configure Domain and SSL

### Option A: Cloudflare (Recommended)
1. **Add domain to Cloudflare**:
   - Sign up at [Cloudflare](https://cloudflare.com)
   - Add your domain and update nameservers
   - Enable "Full (Strict)" SSL mode

2. **Generate Origin Certificates**:
   - Go to SSL/TLS → Origin Server
   - Click "Create Certificate"
   - Select your domain and subdomains
   - Choose "PEM" format
   - Download both `.pem` and `.key` files

3. **Install certificates on server**:
   ```bash
   sudo mkdir -p /etc/nginx/ssl
   sudo cp yourdomain.pem /etc/nginx/ssl/yourdomain.pem
   sudo cp yourdomain.key /etc/nginx/ssl/yourdomain.key
   sudo chmod 600 /etc/nginx/ssl/*
   sudo chown root:root /etc/nginx/ssl/*
   ```

### Option B: Let's Encrypt (Alternative)
```bash
# Install Certbot
sudo apt update
sudo apt install certbot python3-certbot-nginx

# Stop any running nginx
sudo systemctl stop nginx

# Generate certificates
sudo certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com

# Certificates will be in:
# /etc/letsencrypt/live/yourdomain.com/
```

---

## 4. Configure Nginx

> [!NOTE] 
>  Domain and certificates names will be updated automatically, if you decide to use our existing GitHub Actions workflow. But you need to set the Github Secrets first (see section 13).

The nginx configuration will be automatically generated from `backend/nginx-container/nginx.conf.template` using your domain settings.

---

## 5. Setup External Volumes (Production Safety)
Create external Docker volumes for persistent data:
```bash
# Create external volumes for critical data
docker volume create fluently_pgdata_safe
docker volume create fluently_grafana_data_external
docker volume create fluently_prometheus_data_external

# Verify volumes were created
docker volume ls | grep fluently
```

---

## 6. Configure API Keys

### LLM API Keys
1. **Groq API Keys**:
   - Visit [Groq Console](https://console.groq.com/keys)
   - Create new API key
   - Add to `.env` file: `GROQ_API_KEYS=your_key_here`

2. **Gemini API Keys**:
   - Visit [Google MakerSuite](https://makersuite.google.com/app/apikey)
   - Create new API key
   - Add to `.env` file: `GEMINI_API_KEYS=your_key_here`

### Google OAuth Setup
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project or select existing
3. Enable Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URIs:
   - `https://yourdomain.com/auth/google/callback`
   - `https://yourdomain.com/api/v1/auth/google/callback`

### Telegram Bot Setup
1. Message [@BotFather](https://t.me/BotFather) on Telegram
2. Create new bot with `/newbot`
3. Copy bot token to `.env` file
4. Set webhook URL: `https://yourdomain.com/webhook`

---

## 7. Deploy Services
```bash
# Build and start all services
docker compose up -d --build

# Check service status
docker compose ps

# View logs for specific service
docker compose logs -f backend
docker compose logs -f llm-api
docker compose logs -f telegram-bot
```

---

## 8. Service Architecture Overview

The Fluently platform consists of the following services:

### Core Services
- **Backend API** (Port 8070): Main Go application with REST API
- **Postgres Database** (Port 5432): Primary data storage
- **Redis** (Port 6379): Session storage and message queuing
- **Nginx** (Ports 80/443): Reverse proxy and SSL termination

### AI/ML Services
- **LLM API** (Port 8003): AI conversation service with Groq/Gemini
- **ML API** (Port 8001): Machine learning distractor generation
- **Thesaurus API** (Port 8002): Vocabulary recommendations

### Communication Services
- **Telegram Bot** (Port 8060): Telegram bot webhook handler

### Administrative Services
- **Directus CMS** (Port 8055): Content management system
- **Grafana** (Port 3000): Monitoring dashboards
- **Prometheus** (Port 9090): Metrics collection
- **Loki** (Port 3100): Log aggregation

---

## 9. Access Services

### Public Services (via Domain)
- **Main Website**: `https://yourdomain.com/`
- **Backend API**: `https://yourdomain.com/api/v1/`
- **API Documentation**: `https://yourdomain.com/swagger/`
- **Telegram Bot Webhook**: `https://yourdomain.com/webhook`

### Administrative Services (ZeroTier/SSH Tunnel)
- **Directus CMS**: `http://SERVER_IP:8055`
- **Grafana**: `http://SERVER_IP:3000`
- **Prometheus**: `http://SERVER_IP:9090`

### API Service Health Checks
- **Backend Health**: `https://yourdomain.com/health`
- **LLM API Health**: `http://SERVER_IP:8003/health`
- **ML API Health**: `http://SERVER_IP:8001/health`
- **Thesaurus API Health**: `http://SERVER_IP:8002/health`

---

## 10. Security Configuration

### Firewall Setup
```bash
# Allow only necessary ports
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp     # SSH
sudo ufw allow 80/tcp     # HTTP
sudo ufw allow 443/tcp    # HTTPS
sudo ufw enable

# Verify firewall status
sudo ufw status verbose
```

### SSH Security
```bash
# Disable password authentication
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart sshd
```

### ZeroTier VPN Setup (Recommended for Admin Access)
```bash
# Install ZeroTier
curl -s https://install.zerotier.com | sudo bash

# Join your network
sudo zerotier-cli join YOUR_NETWORK_ID

# Verify connection
sudo zerotier-cli listnetworks
```

### Environment Security
```bash
# Secure .env file
chmod 600 .env
sudo chown root:root .env
```

---

## 11. Monitoring Setup

### Configure Grafana
1. Access Grafana at `http://YOUR_ZEROTIER_IP:3000`
2. Login with credentials from `.env` file:
   - Username: `admin`
   - Password: `${GRAFANA_ADMIN_PASSWORD}`
3. Add data sources:
   - **Prometheus**: `http://prometheus:9090`
   - **Loki**: `http://loki:3100`
4. Import dashboards from `backend/monitoring/grafana/dashboards/`

### Health Check Commands
```bash
# Check all services status
docker compose ps

# Verify service health
curl -f https://yourdomain.com/health
curl -f http://localhost:8070/health
curl -f http://localhost:8001/health
curl -f http://localhost:8002/health
curl -f http://localhost:8003/health

# Check LLM API configuration
curl -f http://localhost:8003/config
```

---

## 12. Backup Configuration

### Setup Automated Backups
```bash
# Create backup directory
sudo mkdir -p /home/backups
sudo chown $USER:$USER /home/backups

# Create backup script
cat > /home/backups/backup_fluently.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/home/backups"
DATE=$(date +%Y%m%d_%H%M%S)
PROJECT_DIR="/path/to/Fluently-fork"

cd $PROJECT_DIR

# Database backup
docker compose exec -T postgres pg_dump -U postgres postgres | gzip > "$BACKUP_DIR/postgres_backup_$DATE.sql.gz"

# Volume backups
docker run --rm -v fluently_pgdata_safe:/data -v $BACKUP_DIR:/backup ubuntu tar czf /backup/pgdata_backup_$DATE.tar.gz /data
docker run --rm -v fluently_grafana_data_external:/data -v $BACKUP_DIR:/backup ubuntu tar czf /backup/grafana_backup_$DATE.tar.gz /data

# Configuration backup
tar czf "$BACKUP_DIR/config_backup_$DATE.tar.gz" .env docker-compose.yml

# Cleanup old backups (keep last 7 days)
find $BACKUP_DIR -name "*backup_*.gz" -type f -mtime +7 -delete

echo "Backup completed: $DATE"
EOF

chmod +x /home/backups/backup_fluently.sh
```

### Setup Cron Job
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /home/backups/backup_fluently.sh >> /var/log/fluently_backup.log 2>&1
```

---

## 13. CI/CD Integration (Optional)

Configure GitHub Actions for automated deployment. Add these secrets to your repository:

```yaml
# GitHub Repository → Settings → Secrets and variables → Actions
DEPLOY_HOST: your-server-ip
DEPLOY_USERNAME: your-deploy-username
DEPLOY_SSH_KEY: |
  -----BEGIN OPENSSH PRIVATE KEY-----
  your-private-key-content
  -----END OPENSSH PRIVATE KEY-----
ZEROTIER_IP: your-zerotier-ip
GROQ_API_KEYS: your-groq-api-keys
GEMINI_API_KEYS: your-gemini-api-keys
DB_PASSWORD: your-database-password
JWT_SECRET: your-jwt-secret
DIRECTUS_ADMIN_PASSWORD: your-directus-password
GRAFANA_ADMIN_PASSWORD: your-grafana-password
```

The CI/CD pipeline is configured in `.github/workflows/deploy.yml` and will:
- Build and push Docker images
- Deploy to your server
- Run health checks
- Send notifications

---

## 14. Troubleshooting

### Common Issues

**1. Docker Build Failures**
```bash
# Clear build cache
docker builder prune -a

# Rebuild specific service
docker compose build --no-cache llm-api

# Check service logs
docker compose logs -f llm-api
```

**2. SSL Certificate Issues**
```bash
# Verify certificates
sudo openssl x509 -in /etc/nginx/ssl/yourdomain.pem -text -noout
sudo openssl rsa -in /etc/nginx/ssl/yourdomain.key -check

# Test SSL connection
openssl s_client -connect yourdomain.com:443
```

**3. Database Connection Issues**
```bash
# Check database status
docker compose exec postgres pg_isready -U postgres

# Connect to database
docker compose exec postgres psql -U postgres -d postgres

# View database logs
docker compose logs postgres
```

**4. LLM API Issues**
```bash
# Check API configuration
curl http://localhost:8003/config

# Verify API keys
docker compose exec llm-api env | grep -E "(GROQ|GEMINI)_API_KEYS"

# Test API functionality
curl -X POST http://localhost:8003/chat \
  -H "Content-Type: application/json" \
  -d '{"messages": [{"role": "user", "content": "Hello"}], "model_type": "fast"}'
```

**5. Port Binding Issues**
```bash
# Check port usage
sudo netstat -tulpn | grep :8070
sudo lsof -i :8070

# Change port in .env if needed
APP_PORT=8071
```

---

## 15. Maintenance Commands

### Service Management
```bash
# Stop all services
docker compose down

# Restart specific service
docker compose restart backend

# Update and redeploy
git pull origin main
docker compose up -d --build

# View service logs
docker compose logs -f --tail=100 backend
```

### Database Maintenance
```bash
# Database backup
docker compose exec postgres pg_dump -U postgres postgres > backup.sql

# Database restore
docker compose exec -T postgres psql -U postgres -d postgres < backup.sql

# Check database connections
docker compose exec postgres psql -U postgres -c "SELECT * FROM pg_stat_activity;"
```

### System Cleanup
```bash
# Remove unused Docker resources
docker system prune -a

# Remove old images
docker image prune -f

# Remove unused volumes (be careful!)
docker volume prune

# Check disk usage
docker system df
```

---

## 16. Performance Optimization

### Database Optimization
```bash
# Monitor database performance
docker compose exec postgres psql -U postgres -c "
SELECT query, calls, total_time, mean_time 
FROM pg_stat_statements 
ORDER BY total_time DESC LIMIT 10;"
```

### Resource Monitoring
```bash
# Monitor container resources
docker stats

# Check system resources
htop
df -h
free -h
```

### Service Scaling
```bash
# Scale specific services
docker compose up -d --scale llm-api=2
docker compose up -d --scale thesaurus-api=2
```

---

## Notes

### Important Considerations
- **Resource Requirements**: Minimum 8GB RAM, 4 CPU cores recommended
- **Storage**: At least 50GB for production deployment
- **Network**: Stable internet connection for AI API services
- **Monitoring**: Regularly check service health and logs
- **Security**: Keep all secrets secure and rotate them regularly
- **Backups**: Implement automated backup strategy
- **Updates**: Regularly update dependencies and security patches

### Production Checklist
- [ ] Domain configured with SSL certificates
- [ ] All API keys configured and tested
- [ ] External volumes created for data persistence
- [ ] Firewall configured with minimum required ports
- [ ] SSH hardened with key-based authentication
- [ ] ZeroTier VPN configured for admin access
- [ ] Monitoring dashboards set up in Grafana
- [ ] Backup automation configured and tested
- [ ] CI/CD pipeline configured (optional)
- [ ] All service health checks passing

For local development and testing, see the [Local Installation Guide](Install_Local.md).

---

## Support

If you encounter issues:
1. Check the [troubleshooting section](#14-troubleshooting)
2. Review service logs: `docker compose logs -f [service-name]`
3. Verify configuration: Check `.env` file and service health endpoints
4. Consult the [GitHub repository](https://github.com/FluentlyOrg/Fluently-fork) for updates and issues
