# Fluently Services Quick Reference Card

## 🌐 Main Application
- **Production**: https://fluently-app.ru
- **Staging**: https://fluently-app.online

## 🔐 ZeroTier IPs
- **Production**: `10.243.92.227`
- **Staging**: `10.243.191.108`

## 📊 Monitoring Stack

### Grafana Dashboards
- **Production**: http://10.243.92.227:3000
- **Staging**: http://10.243.191.108:3000
- **Login**: `admin` / `admin123` (check `.env` for current password)

### Prometheus Metrics
- **Production**: http://10.243.92.227:9090
- **Staging**: http://10.243.191.108:9090

### Loki Logs
- **Production**: http://10.243.92.227:3100
- **Staging**: http://10.243.191.108:3100

## 🗄️ Database & CMS

### Directus CMS
- **Production**: http://10.243.92.227:8055
- **Staging**: http://10.243.191.108:8055

### PostgreSQL Database
- **Production**: `10.243.92.227:5432`
- **Staging**: `10.243.191.108:5432`
- **Database**: `postgres`
- **User**: Check `.env` file for credentials

## 🤖 ML Services

### ML API (Distractor Generation)
- **Production**: http://10.243.92.227:8001
- **Staging**: http://10.243.191.108:8001
- **Health Check**: `/health`

## 🔍 Code Quality

### SonarQube
- **Production**: http://10.243.92.227:9000
- **Staging**: http://10.243.191.108:9000

## 📈 Metric Exporters

| Service | Production | Staging | Purpose |
|---------|------------|---------|---------|
| Node Exporter | 10.243.92.227:9100 | 10.243.191.108:9100 | System metrics |
| Nginx Exporter | 10.243.92.227:9113 | 10.243.191.108:9113 | Web server metrics |
| cAdvisor | 10.243.92.227:8044 | 10.243.191.108:8044 | Container metrics |

## 🔧 SSH Access
```bash
# Production
ssh deploy@fluently-app.ru

# Staging
ssh deploy-staging@fluently-app.online
```

## 🌐 SSH Port Forwarding (Alternative to ZeroTier)

### Forward All Services - Staging
```bash
ssh -L 3000:localhost:3000 \
    -L 8055:localhost:8055 \
    -L 9000:localhost:9000 \
    -L 9090:localhost:9090 \
    -L 5432:localhost:5432 \
    -L 8070:localhost:8070 \
    -L 8001:localhost:8001 \
    deploy-staging@fluently-app.online
```

### Forward All Services - Production
```bash
ssh -L 3000:localhost:3000 \
    -L 8055:localhost:8055 \
    -L 9000:localhost:9000 \
    -L 9090:localhost:9090 \
    -L 5432:localhost:5432 \
    -L 8070:localhost:8070 \
    -L 8001:localhost:8001 \
    deploy@fluently-app.ru
```

### Port Reference
- **3000**: Grafana
- **8055**: Directus CMS
- **9000**: SonarQube
- **9090**: Prometheus
- **5432**: PostgreSQL
- **8070**: Backend API
- **8001**: ML API

## ⚡ Quick Commands
```bash
# Service status
docker compose ps

# View logs
docker compose logs -f [service]

# Restart service
docker compose restart [service]

# Health checks
curl http://localhost:8070/health  # Backend
curl http://localhost:8001/health  # ML API

# Database access
docker compose exec postgres psql -U postgres -d postgres

# Check ZeroTier
sudo zerotier-cli status
sudo zerotier-cli listnetworks
```

## 🔄 Local Development
```bash
# Start all services locally
make run-production

# Start core services only
make run-core

# Stop all services
make stop

# Clean up (removes volumes)
make clean
```

## 🆘 Emergency Contacts
- **DevOps**: Timofey - [Contact Info]
- **Backend Lead**: [Contact Info]
- **Frontend Lead**: [Contact Info]

## 📋 Service Dependencies
```
Frontend → Backend API → PostgreSQL
Directus → PostgreSQL  
SonarQube → PostgreSQL
ML API → HuggingFace Models
Grafana → Prometheus + Loki
Nginx → Backend + Static Files
All Exporters → Target Services
```

## 🔐 ZeroTier Setup
```bash
# Install ZeroTier
curl -s https://install.zerotier.com | sudo bash

# Join network (get ID from admin)
sudo zerotier-cli join [NETWORK_ID]

# Check status
sudo zerotier-cli status
```

---
*💡 Tip: Bookmark this page and keep it open during development work*
