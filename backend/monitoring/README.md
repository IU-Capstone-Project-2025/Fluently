# Fluently Monitoring Stack

This directory contains the configuration files for the complete monitoring stack for the Fluently application.

## Services Overview

The monitoring stack includes:

- **Prometheus** - Metrics collection and storage
- **Grafana** - Visualization and dashboards  
- **Loki** - Log aggregation and storage
- **Promtail** - Log shipping agent
- **SonarQube** - Code quality analysis
- **Node Exporter** - System metrics collection

All services are accessible only via the ZeroTier internal IP: `10.243.92.227`

## Service URLs

When connected to ZeroTier network:

- **Grafana**: http://10.243.92.227:3000
  - Username: `admin`
  - Password: Set in `GRAFANA_ADMIN_PASSWORD` env variable

- **Prometheus**: http://10.243.92.227:9090
- **Loki**: http://10.243.92.227:3100
- **SonarQube**: http://10.243.92.227:9000
  - Initial setup required on first visit

## Configuration Files

### Prometheus (`monitoring/prometheus/`)
- `prometheus.yml` - Main Prometheus configuration
- `alert_rules.yml` - Alerting rules for infrastructure monitoring

### Grafana (`monitoring/grafana/`)
- `provisioning/datasources/datasources.yml` - Auto-configured data sources
- `provisioning/dashboards/dashboards.yml` - Dashboard provisioning config
- `dashboards/` - Pre-built dashboards:
  - `infrastructure-overview.json` - System metrics overview
  - `logs-dashboard.json` - Centralized log viewing
  - `api-dashboard.json` - API-specific metrics
  - `system-dashboard.json` - System-level monitoring

### Loki (`monitoring/loki/`)
- `loki-config.yml` - Loki configuration for log storage

### Promtail (`monitoring/promtail/`)
- `promtail-config.yml` - Log collection configuration for:
  - Docker container logs
  - Nginx access/error logs
  - System logs
  - Application logs

### SonarQube (`sonarqube/`)
- `sonar.properties` - SonarQube configuration

## Setup Instructions

1. **Environment Variables**: Copy `.env.example` to `.env` and configure:
   ```bash
   cp .env.example .env
   # Edit .env with your specific values
   ```

2. **Start the Stack**:
   ```bash
   docker-compose up -d
   ```

3. **Access via ZeroTier**:
   - Ensure you're connected to the ZeroTier network
   - Access services using the internal IP: `10.243.92.227`

4. **Initial SonarQube Setup**:
   - Visit http://10.243.92.227:9000
   - Login with default credentials (admin/admin)
   - Change the default password
   - Create a new project for your Go application

## Security Features

- **Network Isolation**: All admin services bound to ZeroTier IP only
- **No Public Exposure**: Monitoring tools are not accessible from the internet
- **Authentication**: All services require authentication
- **Encrypted Communication**: ZeroTier provides encrypted network tunnel

## Monitoring Features

### Metrics Collected
- System metrics (CPU, memory, disk, network)
- Application metrics (API response times, error rates)
- Container metrics (Docker stats)
- Database metrics (PostgreSQL performance)
- Web server metrics (Nginx stats)

### Logs Collected
- Application logs
- System logs (/var/log/syslog)
- Nginx access and error logs
- Docker container logs
- All structured with proper labels for easy filtering

### Alerts Configured
- High CPU usage (>80%)
- High memory usage (>85%)
- Service downtime
- Low disk space (<15%)
- Application errors
- Slow response times

## Usage Examples

### View Logs in Grafana
1. Go to Grafana → Explore
2. Select "Loki" datasource
3. Use queries like:
   - `{job="fluently-app"}` - Application logs
   - `{job="nginx-access"}` - Nginx access logs
   - `{container_name="fluently_app"}` - Specific container logs

### Create Custom Dashboards
1. Go to Grafana → Dashboards → New
2. Add panels with Prometheus queries:
   - `rate(http_requests_total[5m])` - Request rate
   - `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))` - 95th percentile response time

### SonarQube Analysis
1. In your project directory:
   ```bash
   # Install SonarScanner
   # Run analysis
   sonar-scanner \
     -Dsonar.projectKey=fluently-app \
     -Dsonar.sources=. \
     -Dsonar.host.url=http://10.243.92.227:9000 \
     -Dsonar.login=your_token
   ```

## Troubleshooting

### Services Not Accessible
- Check ZeroTier connection: `zerotier-cli info`
- Verify containers are running: `docker-compose ps`
- Check port bindings: `netstat -tlnp | grep 10.243.92.227`

### High Resource Usage
- Monitor resource usage in Grafana
- Adjust memory limits in docker-compose.yml if needed
- Consider data retention policies for Prometheus and Loki

### Log Collection Issues
- Check Promtail container logs: `docker-compose logs promtail`
- Verify log file permissions and paths
- Test Loki connectivity from Promtail

## Data Retention

- **Prometheus**: 200 hours (configurable in prometheus.yml)
- **Loki**: Default retention (configurable in loki-config.yml)
- **Grafana**: Persistent storage via Docker volumes

## Backup Recommendations

Regular backups of Docker volumes:
```bash
docker-compose down
sudo tar czf monitoring-backup-$(date +%Y%m%d).tar.gz \
  /var/lib/docker/volumes/backend_prometheus_data \
  /var/lib/docker/volumes/backend_grafana_data \
  /var/lib/docker/volumes/backend_loki_data \
  /var/lib/docker/volumes/backend_sonarqube_data
docker-compose up -d
```
