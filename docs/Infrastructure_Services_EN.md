# Fluently Infrastructure Services

This document describes all available services in the Fluently project infrastructure for both production and staging environments.

## 🌐 Environment Overview

- **Production Domain**: `fluently-app.ru`
- **Staging Domain**: `fluently-app.online`
- **Production ZeroTier IP**: `10.243.92.227`
- **Staging ZeroTier IP**: `10.243.191.108`

> **Note**: All monitoring and development services are accessible only via ZeroTier VPN for security reasons.

## 📱 Main Application Services

### Web Application
- **Production**: https://fluently-app.ru
- **Staging**: https://fluently-app.online
- **Purpose**: Main Fluently language learning web application with user interface and API endpoints.

### API Health Check
- **Production**: https://fluently-app.ru/health
- **Staging**: https://fluently-app.online/health
- **Purpose**: Application health status endpoint for monitoring service availability.

## 🗄️ Database Services

### PostgreSQL Database
- **Production**: `10.243.92.227:5432`
- **Staging**: `10.243.191.108:5432`
- **Purpose**: Primary database storing user data, lessons, vocabulary, and application state.

### Directus CMS
- **Production**: `http://10.243.92.227:8055`
- **Staging**: `http://10.243.191.108:8055`
- **Purpose**: Headless CMS for content management, database administration, and API generation.

## 📊 Monitoring & Analytics

### Grafana (Dashboards)
- **Production**: `http://10.243.92.227:3000`
- **Staging**: `http://10.243.191.108:3000`
- **Purpose**: Visual monitoring dashboards for system metrics, application performance, and infrastructure health.
- **Default credentials**: `admin` / `admin123` (check `.env` for current password)

### Prometheus (Metrics)
- **Production**: `http://10.243.92.227:9090`
- **Staging**: `http://10.243.191.108:9090`
- **Purpose**: Metrics collection and time-series database for system and application monitoring.

### Loki (Logs)
- **Production**: `http://10.243.92.227:3100`
- **Staging**: `http://10.243.191.108:3100`
- **Purpose**: Centralized log aggregation and search for application and system logs.

## 🔍 Code Quality & Security

### SonarQube
- **Production**: `http://10.243.92.227:9000`
- **Staging**: `http://10.243.191.108:9000`
- **Purpose**: Static code analysis, security vulnerability detection, and code quality metrics.

## 📈 Metric Exporters

### Node Exporter (System Metrics)
- **Production**: `http://10.243.92.227:9100`
- **Staging**: `http://10.243.191.108:9100`
- **Purpose**: Exports system-level metrics like CPU, memory, disk usage, and network statistics.

### PostgreSQL Exporter
- **Production**: `http://10.243.92.227:9187`
- **Staging**: `http://10.243.191.108:9187`
- **Purpose**: Exports PostgreSQL database metrics for monitoring database performance and health.

### Nginx Exporter
- **Production**: `http://10.243.92.227:9113`
- **Staging**: `http://10.243.191.108:9113`
- **Purpose**: Exports Nginx web server metrics including request rates, response times, and connection statistics.

### cAdvisor (Container Metrics)
- **Production**: `http://10.243.92.227:8044`
- **Staging**: `http://10.243.191.108:8044`
- **Purpose**: Exports Docker container metrics including resource usage, performance characteristics, and container lifecycle information.

## 🔐 Access Requirements

### ZeroTier VPN
- **Network ID**: Contact admin for network invitation
- **Purpose**: Secure access to development and monitoring services
- **Required for**: All services except main web applications

### Service Accounts
- **Grafana**: `admin` / (see `.env` file)
- **SonarQube**: Setup required on first access
- **Directus**: Configure admin account during setup

## 🚀 Deployment Information

### Docker Services
All services run in Docker containers managed by Docker Compose:
- **Production**: `/home/deploy/Fluently-fork/backend/`
- **Staging**: `/home/deploy-staging/Fluently-fork/backend/`

### SSL Configuration
- **Production**: Cloudflare Full (Strict) SSL with Origin Certificates
- **Staging**: Cloudflare Full SSL mode
- **Management**: Use `./nginx-container/manage-nginx.sh` script

## 📝 Quick Access Commands

### SSH Access
```bash
# Production
ssh deploy@fluently-app.ru

# Staging  
ssh deploy-staging@fluently-app.online
```

### Service Management
```bash
# View running services
docker compose ps

# View service logs
docker compose logs [service_name]

# Restart services
docker compose restart [service_name]
```

### Monitoring URLs Quick Reference
```bash
# Production Monitoring Stack
http://10.243.92.227:3000  # Grafana
http://10.243.92.227:9090  # Prometheus
http://10.243.92.227:9000  # SonarQube
http://10.243.92.227:8055  # Directus
http://10.243.92.227:9100  # Node Exporter
http://10.243.92.227:9187  # PostgreSQL Exporter
http://10.243.92.227:9113  # Nginx Exporter
http://10.243.92.227:8044  # cAdvisor
http://10.243.92.227:3100  # Loki

# Staging Monitoring Stack
http://10.243.191.108:3000  # Grafana
http://10.243.191.108:9090  # Prometheus
http://10.243.191.108:9000  # SonarQube
http://10.243.191.108:8055  # Directus
http://10.243.191.108:9100  # Node Exporter
http://10.243.191.108:9187  # PostgreSQL Exporter
http://10.243.191.108:9113  # Nginx Exporter
http://10.243.191.108:8044  # cAdvisor
http://10.243.191.108:3100  # Loki
```

## 🆘 Support & Troubleshooting

### Health Checks
1. **Application**: Check `/health` endpoint
2. **Database**: Verify PostgreSQL exporter metrics
3. **Services**: Use `docker compose ps` to check container status

### Log Access
- **Application logs**: `docker compose logs app`
- **System logs**: Available in Grafana via Loki
- **Nginx logs**: `docker compose logs nginx`

### Emergency Contacts
- **DevOps Admin**: Timofey (Lead Infrastructure) - [Add contact]
- **Backend Team Lead**: [Add contact]
- **Frontend Team Lead**: [Add contact]
- **Project Manager**: [Add contact]

### Common Troubleshooting Steps
1. **Service Down**: Check `docker compose ps` and restart if needed
2. **502/503 Errors**: Verify app container is running and healthy
3. **SSL Issues**: Check certificate validity and Nginx config
4. **Performance Issues**: Check Grafana dashboards for resource usage
5. **Database Issues**: Check PostgreSQL logs and connection pool status

### Useful Commands
```bash
# Check all container statuses
docker compose ps

# View application logs
docker compose logs -f app

# Restart all services
docker compose restart

# Check Nginx configuration syntax
docker compose exec nginx nginx -t

# View system resource usage
docker stats

# Check ZeroTier connection status
sudo zerotier-cli status
```

### Port Forwarding for Local Development
If you need to access services locally without ZeroTier:
```bash
# Forward Grafana (use with caution)
ssh -L 3000:10.243.92.227:3000 deploy@fluently-app.ru

# Forward PostgreSQL (for database tools)
ssh -L 5432:10.243.92.227:5432 deploy@fluently-app.ru
```

### Service Dependencies
- **App** depends on: PostgreSQL
- **Directus** depends on: PostgreSQL
- **SonarQube** depends on: PostgreSQL
- **Grafana** depends on: Prometheus, Loki
- **Exporters** depend on: Their target services (Nginx, PostgreSQL, etc.)
- **Nginx** serves: Main application, static files

---

