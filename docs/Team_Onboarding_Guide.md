# Fluently Team Onboarding Guide

Welcome to the Fluently development team! This guide will help you get started with our infrastructure and development workflow.

## 🚀 Quick Start Checklist

### Prerequisites
- [ ] Receive ZeroTier network invitation from DevOps admin
- [ ] Get SSH access credentials for servers
- [ ] Install required tools on your local machine
- [ ] Configure development environment

### Required Tools Installation

#### 1. ZeroTier VPN Client
```bash
# Ubuntu/Debian
curl -s https://install.zerotier.com | sudo bash

# macOS
brew install zerotier-one

# Windows
# Download installer from https://www.zerotier.com/download/
```

#### 2. Docker & Docker Compose
```bash
# Ubuntu/Debian
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# macOS
brew install docker docker-compose

# Add user to docker group (Linux)
sudo usermod -aG docker $USER
```

#### 3. Development Tools
```bash
# Git, SSH, basic tools
sudo apt install git openssh-client curl wget

# Optional: VS Code for remote development
sudo snap install code --classic
```

## 🔐 Access Setup

### 1. ZeroTier Network Access
1. Install ZeroTier client
2. Join network with ID provided by admin:
   ```bash
   sudo zerotier-cli join [NETWORK_ID]
   ```
3. Wait for admin approval
4. Verify connection:
   ```bash
   sudo zerotier-cli status
   sudo zerotier-cli listnetworks
   ```

### 2. SSH Key Configuration
```bash
# Generate SSH key if you don't have one
ssh-keygen -t ed25519 -C "your.email@example.com"

# Copy public key to admin
cat ~/.ssh/id_ed25519.pub
```

### 3. Test Server Access
```bash
# Production server
ssh deploy@fluently-app.ru

# Staging server
ssh deploy-staging@fluently-app.online
```

## 🌐 Service Access Guide

### Main Application URLs
- **Production**: https://fluently-app.ru
- **Staging**: https://fluently-app.online

### Development & Monitoring (ZeroTier Required)

#### Grafana Dashboards
- **URL**: http://10.243.92.227:3000 (Production) / http://10.243.191.108:3000 (Staging)
- **Default Login**: `admin` / `admin123` (check `.env` for current password)
- **Purpose**: Visual monitoring, metrics, alerts

#### Directus CMS
- **URL**: http://10.243.92.227:8055 (Production) / http://10.243.191.108:8055 (Staging)
- **Purpose**: Content management, database administration
- **Setup**: Create admin account on first visit

#### SonarQube Code Quality
- **URL**: http://10.243.92.227:9000 (Production) / http://10.243.191.108:9000 (Staging)
- **Purpose**: Code analysis, security scans, quality metrics
- **Setup**: Configure on first visit

#### Prometheus Metrics
- **URL**: http://10.243.92.227:9090 (Production) / http://10.243.191.108:9090 (Staging)
- **Purpose**: Raw metrics data, alerting rules

## 🛠️ Development Workflow

### 1. Working with Docker Services
```bash
# Navigate to backend directory
cd /path/to/Fluently-fork/backend

# Check running services
docker compose ps

# View logs
docker compose logs [service_name]

# Restart services
docker compose restart [service_name]

# Rebuild and restart
docker compose up --build -d [service_name]
```

### 2. Database Access
```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U [username] -d [database]

# Or via Directus web interface
# http://10.243.92.227:8055
```

### 3. Log Analysis
```bash
# Application logs
docker compose logs -f app

# Nginx logs
docker compose logs -f nginx

# Database logs
docker compose logs -f postgres

# All logs with timestamps
docker compose logs -f --timestamps
```

## 📊 Monitoring & Debugging

### Common Issues & Solutions

#### Service Not Responding
1. Check container status: `docker compose ps`
2. Check logs: `docker compose logs [service]`
3. Restart if needed: `docker compose restart [service]`

#### 502/503 Errors
1. Verify app container is running
2. Check application logs for errors
3. Verify database connectivity
4. Check Nginx configuration

#### Database Connection Issues
1. Check PostgreSQL container status
2. Verify database credentials in `.env`
3. Check connection pool status in logs
4. Restart PostgreSQL if needed

#### SSL/Certificate Issues
1. Check certificate expiration dates
2. Verify Nginx SSL configuration
3. Check Cloudflare SSL settings
4. Review certificate installation

### Grafana Dashboard Usage
1. **System Overview**: CPU, memory, disk usage
2. **Application Metrics**: Request rates, response times, errors
3. **Database Performance**: Query performance, connections, locks
4. **Infrastructure Health**: Service availability, resource alerts

### Using Prometheus for Debugging
- Query metrics directly: `http://10.243.92.227:9090`
- Example queries:
  - CPU usage: `100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
  - Memory usage: `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`
  - HTTP requests: `rate(nginx_http_requests_total[5m])`

## 🔧 Local Development Setup

### Environment Variables
Create `.env` file in backend directory:
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fluently
DB_USER=your_username
DB_PASSWORD=your_password

# Application
PORT=8070
NODE_ENV=development

# Directus
DIRECTUS_SECRET_KEY=your_secret_key
DIRECTUS_ADMIN_EMAIL=admin@example.com
DIRECTUS_ADMIN_PASSWORD=admin_password

# Monitoring
GRAFANA_ADMIN_PASSWORD=secure_password
```

### Running Services Locally
```bash
# Start only necessary services for development
docker compose up postgres directus -d

# Run application locally
go run cmd/main.go

# Or start all services
docker compose up -d
```

## 📝 Code Contribution Guidelines

### 1. Git Workflow
```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Make changes and commit
git add .
git commit -m "feat: add new feature"

# Push and create pull request
git push origin feature/your-feature-name
```

### 2. Code Quality Checks
- Run SonarQube analysis before pushing
- Ensure all tests pass
- Check logs for any errors
- Verify monitoring dashboards show healthy metrics

### 3. Deployment Process
- Changes are automatically deployed via GitHub Actions
- Monitor deployment in Actions tab
- Check application health after deployment
- Verify monitoring dashboards post-deployment

## 📞 Support & Contacts

...


## 🎯 Next Steps After Onboarding

1. [ ] Complete all access setup tasks
2. [ ] Explore Grafana dashboards
3. [ ] Familiarize yourself with codebase
4. [ ] Set up local development environment
5. [ ] Create test pull request
6. [ ] Join team communication channels
7. [ ] Schedule introduction meeting with team lead

---