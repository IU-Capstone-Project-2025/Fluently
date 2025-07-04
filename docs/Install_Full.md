# Full Production Installation Guide

This guide explains how to deploy Fluently in a production-like environment with your own domain, SSL certificates, and all services running as in production.

---

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)
- Your own domain name (e.g., `your-domain.com`)
- VPS or cloud server with public IP
- Basic understanding of DNS configuration

---

## 1. Clone the Repository
```bash
git clone https://github.com/IU-Capstone-Project-2025/Fluently.git
cd Fluently-fork
```

---

## 2. Create Environment File
Create a `.env` file in the `backend` directory:
```env
# JWT
JWT_SECRET=supersecretjwtkey_CHANGE_THIS
JWT_EXPIRATION=24h
REFRESH_EXPIRATION=720h

# App
APP_NAME=fluently
APP_PORT=8070

# Database
DB_USER=postgres
DB_PASSWORD=secure_password_CHANGE_THIS
DB_HOST=postgres
DB_PORT=5432
DB_NAME=postgres

# Directus configuration
DIRECTUS_PORT=8055
DIRECTUS_ADMIN_EMAIL=admin@yourdomain.com
DIRECTUS_ADMIN_PASSWORD=admin_password_CHANGE_THIS
DIRECTUS_SECRET_KEY=directus_secret_CHANGE_THIS

# Grafana
GRAFANA_ADMIN_PASSWORD=grafana_password_CHANGE_THIS

# Google OAuth (Configure in Google Cloud Console)
IOS_GOOGLE_CLIENT_ID=your-ios-client-id.apps.googleusercontent.com
ANDROID_GOOGLE_CLIENT_ID=your-android-client-id.apps.googleusercontent.com
WEB_GOOGLE_CLIENT_ID=your-web-client-id.apps.googleusercontent.com

# Server configuration
ZEROTIER_IP=YOUR_SERVER_IP
PUBLIC_URL=https://your-domain.com
```

Copy to project root for Docker Compose:
```bash
cp backend/.env .env
```

---

## 3. Configure Domain and SSL

### Option A: Cloudflare (Recommended)
1. **Add domain to Cloudflare**
2. **Generate Origin Certificates:**
   - Go to SSL/TLS → Origin Server
   - Create Certificate for your domain
   - Download `.pem` and `.key` files

3. **Install certificates:**
   ```bash
   sudo mkdir -p /etc/nginx/ssl
   sudo cp your-domain.pem /etc/nginx/ssl/your-domain.pem
   sudo cp your-domain.key /etc/nginx/ssl/your-domain.key
   sudo chmod 600 /etc/nginx/ssl/*
   ```

4. **Update nginx template:**
   ```bash
   # Edit backend/nginx-container/nginx.conf.template
   # Replace ${DOMAIN} and ${CERT_NAME} placeholders
   ```

### Option B: Let's Encrypt
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Generate certificates
sudo certbot certonly --nginx -d your-domain.com -d www.your-domain.com

# Certificates will be in:
# /etc/letsencrypt/live/your-domain.com/
```

---

## 4. Configure Nginx
Edit `backend/nginx-container/nginx.conf.template`:
```nginx
server_name your-domain.com www.your-domain.com;

# Update certificate paths
ssl_certificate /etc/nginx/ssl/your-domain.pem;
ssl_certificate_key /etc/nginx/ssl/your-domain.key;
```

---

## 5. Setup External Volumes (Production Safety)
```bash
# Create external volumes for data persistence
docker volume create fluently_pgdata_safe
docker volume create fluently_grafana_data_external
docker volume create fluently_prometheus_data_external
docker volume create fluently_sonarqube_data_external
```

---

## 6. Start All Services
```bash
# Generate nginx config with your domain
export DOMAIN="your-domain.com"
export CERT_NAME="your-domain"
envsubst '${DOMAIN} ${CERT_NAME}' < backend/nginx-container/nginx.conf.template > backend/nginx-container/default.conf

# Start services
docker compose up -d --build
```

---

## 7. Access Services

### Public Services
- **Backend API:** `https://your-domain.com/api/v1/`
- **Swagger UI:** `https://your-domain.com/swagger/`
- **Frontend:** `https://your-domain.com/`

### Administrative Services (Restrict Access)
- **Directus CMS:** `https://your-domain.com/admin/` or `http://SERVER_IP:8055`
- **Grafana:** `http://SERVER_IP:3000`
- **Prometheus:** `http://SERVER_IP:9090`
- **SonarQube:** `http://SERVER_IP:9000`

> **Security Note:** Administrative services should be accessed via VPN or SSH tunneling in production.

---

## 8. Security Configuration

### Firewall Setup
```bash
# Allow only necessary ports
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw enable

# Block direct access to admin ports from internet
sudo ufw deny 3000   # Grafana
sudo ufw deny 8055   # Directus
sudo ufw deny 9000   # SonarQube
sudo ufw deny 9090   # Prometheus
```

### SSH Hardening
```bash
# Disable password authentication
sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
sudo systemctl restart sshd
```

### ZeroTier VPN (Optional but Recommended)
```bash
# Install ZeroTier
curl -s https://install.zerotier.com | sudo bash

# Join your network
sudo zerotier-cli join YOUR_NETWORK_ID
```

---

## 9. Monitoring Setup

### Configure Grafana
1. Access Grafana at `http://SERVER_IP:3000`
2. Login with admin credentials from `.env`
3. Import dashboards from `backend/monitoring/grafana/dashboards/`
4. Configure data sources (Prometheus, Loki)

### Health Checks
```bash
# Check all services
docker compose ps

# Test endpoints
curl https://your-domain.com/health
curl http://localhost:8070/health
curl http://localhost:8001/health
```

---

## 10. Backup Configuration

### Setup Automated Backups
```bash
# Create backup directory
mkdir -p /home/backups

# Setup cron for regular backups
crontab -e
# Add: 0 2 * * * /path/to/Fluently-fork/scripts/backup_volumes.sh
```

### Manual Backup
```bash
# Backup volumes
docker compose exec postgres pg_dump -U postgres postgres > backup.sql

# Backup entire project
tar -czf fluently-backup-$(date +%Y%m%d).tar.gz --exclude=node_modules --exclude=.git .
```

---

## 11. CI/CD Integration (Optional)

Configure GitHub Actions with these secrets:
```yaml
DEPLOY_HOST: your-server-ip
DEPLOY_USERNAME: deploy-user
DEPLOY_SSH_KEY: |
  -----BEGIN OPENSSH PRIVATE KEY-----
  your-private-key-content
  -----END OPENSSH PRIVATE KEY-----
ZEROTIER_IP: your-server-ip
```

See [deploy.yml](../.github/workflows/deploy.yml) for the complete workflow.

---

## 12. Maintenance Commands

```bash
# View logs
docker compose logs -f [service_name]

# Restart service
docker compose restart [service_name]

# Update application
git pull origin main
docker compose up -d --build

# Clean up old images
docker image prune -f

# Database maintenance
docker compose exec postgres pg_isready -U postgres
```

---

## Notes
- This setup is intended for production environments
- Requires a real domain and public server
- Use strong passwords and keep them secure
- Regularly update certificates and dependencies
- Monitor logs for security issues
- For local testing, see [Local Installation Guide](Install_Local.md)
