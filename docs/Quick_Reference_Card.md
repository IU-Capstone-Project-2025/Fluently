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
- **Login**: `admin` / `admin123`

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

## 🔍 Code Quality

### SonarQube
- **Production**: http://10.243.92.227:9000
- **Staging**: http://10.243.191.108:9000

## 📈 Metric Exporters

| Service | Production | Staging | Purpose |
|---------|------------|---------|---------|
| Node Exporter | 10.243.92.227:9100 | 10.243.191.108:9100 | System metrics |
| PostgreSQL Exporter | 10.243.92.227:9187 | 10.243.191.108:9187 | Database metrics |
| Nginx Exporter | 10.243.92.227:9113 | 10.243.191.108:9113 | Web server metrics |
| cAdvisor | 10.243.92.227:8044 | 10.243.191.108:8044 | Container metrics |

## 🔧 SSH Access
```bash
# Production
ssh deploy@fluently-app.ru

# Staging
ssh deploy-staging@fluently-app.online
```

## ⚡ Quick Commands
```bash
# Service status
docker compose ps

# View logs
docker compose logs -f [service]

# Restart service
docker compose restart [service]

# Check ZeroTier
sudo zerotier-cli status
```

## 🆘 Emergency Contacts
- **DevOps**: Timofey - [Add contact]
- **Backend Lead**: [Add contact]
- **Frontend Lead**: [Add contact]

## 📋 Service Dependencies
- App → PostgreSQL
- Directus → PostgreSQL  
- SonarQube → PostgreSQL
- Grafana → Prometheus + Loki
- All Exporters → Target Services

---
*Bookmark this page: Keep it handy for daily development work*
