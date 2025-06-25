# Dual Environment Deployment Setup

This document explains how to set up and use the dual-environment deployment system for the Fluently project.

## Overview

The deployment system automatically deploys:
- **Production Environment**: `main` branch → `fluently-app.ru`
- **Staging Environment**: All other branches → `fluently-app.online`

## Required GitHub Secrets

You need to set up the following secrets in your GitHub repository (`Settings > Secrets and variables > Actions`):

### Production Environment
- `PROD_DEPLOY_HOST` - IP address of your production server
- `PROD_DEPLOY_USERNAME` - SSH username for production server
- `PROD_DEPLOY_SSH_KEY` - SSH private key for production server

### Staging Environment
- `STAGING_DEPLOY_HOST` - IP address of your staging server
- `STAGING_DEPLOY_USERNAME` - SSH username for staging server
- `STAGING_DEPLOY_SSH_KEY` - SSH private key for staging server

## Server Setup

### 1. Production Server (fluently-app.ru)
This is your existing server setup. No changes needed.

### 2. Staging Server (fluently-app.online)

#### Server Requirements
- Same specifications as production (or smaller for cost savings)
- Ubuntu/Debian Linux
- Docker and Docker Compose installed
- Git installed
- ZeroTier client (if using ZeroTier monitoring)

#### Initial Setup
```bash
# 1. Create deploy user
sudo adduser deploy
sudo usermod -aG docker deploy
sudo usermod -aG sudo deploy

# 2. Set up SSH key authentication
sudo mkdir -p /home/deploy/.ssh
sudo nano /home/deploy/.ssh/authorized_keys
# (paste your public key)
sudo chown -R deploy:deploy /home/deploy/.ssh
sudo chmod 700 /home/deploy/.ssh
sudo chmod 600 /home/deploy/.ssh/authorized_keys

# 3. Clone repository
sudo -u deploy git clone https://github.com/your-username/Fluently-fork.git /home/deploy/Fluently-fork

# 4. Create backup directory
sudo mkdir -p /home/deploy/backups
sudo chown deploy:deploy /home/deploy/backups
```

## DNS Configuration

### For fluently-app.online:
Add these DNS records:
```
A     fluently-app.online     → [STAGING_SERVER_IP]
A     www.fluently-app.online → [STAGING_SERVER_IP]
A     admin.fluently-app.online → [STAGING_SERVER_IP]
```

### For fluently-app.ru (existing):
```
A     fluently-app.ru         → [PRODUCTION_SERVER_IP]
A     www.fluently-app.ru     → [PRODUCTION_SERVER_IP]
A     admin.fluently-app.ru   → [PRODUCTION_SERVER_IP]
```

## SSL Certificates Setup

### Option 1: Let's Encrypt (Recommended)
```bash
# On staging server
sudo apt install certbot
sudo certbot certonly --standalone -d fluently-app.online -d www.fluently-app.online -d admin.fluently-app.online

# Set up auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Option 2: Cloudflare (if using Cloudflare)
Use Cloudflare's SSL/TLS encryption and origin certificates.

## Environment Configuration

The deployment script automatically updates these environment variables based on the target:

### Production (.env)
```bash
PUBLIC_URL=https://fluently-app.ru
```

### Staging (.env)
```bash
PUBLIC_URL=https://fluently-app.online
```

## Deployment Workflow

### Automatic Deployments
1. **Push to `main`** → Deploys to production (`fluently-app.ru`)
2. **Push to any other branch** → Deploys to staging (`fluently-app.online`)

### Manual Deployments
Go to `Actions > Deploy > Run workflow` and choose:
- Branch to deploy
- Environment (production/staging/auto)

## Monitoring Setup

If using ZeroTier monitoring, configure both servers with the same monitoring stack:

### Production Monitoring URLs (ZeroTier only)
- Grafana: `http://10.243.92.227:3000`
- Prometheus: `http://10.243.92.227:9090`
- SonarQube: `http://10.243.92.227:9000`

### Staging Monitoring URLs (ZeroTier only)
- Grafana: `http://[STAGING_ZEROTIER_IP]:3000`
- Prometheus: `http://[STAGING_ZEROTIER_IP]:9090`
- SonarQube: `http://[STAGING_ZEROTIER_IP]:9000`

## Testing the Setup

### 1. Test Staging Deployment
```bash
git checkout -b feature/test-deployment
git push origin feature/test-deployment
```
Check: https://fluently-app.online

### 2. Test Production Deployment
```bash
git checkout main
git push origin main
```
Check: https://fluently-app.ru

## Troubleshooting

### Common Issues

1. **SSL Certificate Issues**
   ```bash
   # Check certificate status
   sudo certbot certificates
   
   # Renew if needed
   sudo certbot renew
   ```

2. **Docker Permission Issues**
   ```bash
   # Ensure deploy user is in docker group
   sudo usermod -aG docker deploy
   newgrp docker
   ```

3. **Nginx Configuration Issues**
   ```bash
   # Check nginx config
   docker compose exec nginx nginx -t
   
   # View nginx logs
   docker compose logs nginx
   ```

4. **Health Check Failures**
   ```bash
   # Check application logs
   docker compose logs app
   
   # Check if port is accessible
   curl http://localhost:8070/health
   ```

## Backup Strategy

### Production Backups
- Automatic backup before each deployment
- Stored in `/home/deploy/backups/`
- Keep last 7 days of backups

### Staging Backups
- No automatic backups (since it's for testing)
- Manual backup if needed before major changes

## Security Considerations

1. **Firewall Rules**: Both servers should have UFW configured
2. **SSH Access**: Use key-based authentication only
3. **ZeroTier**: Monitoring tools accessible only via ZeroTier
4. **SSL**: Always use HTTPS in production and staging
5. **Environment Variables**: Never commit secrets to git

## Cost Optimization

### Staging Server Options
1. **Smaller VPS**: Use a smaller instance than production
2. **Shared Resources**: Multiple staging environments on one server
3. **Auto-shutdown**: Script to shut down staging after hours (optional)

### Example Auto-shutdown Script (optional)
```bash
# /home/deploy/shutdown-staging.sh
#!/bin/bash
# Shutdown staging environment at night (23:00)
if [ $(date +%H) -eq 23 ]; then
    cd /home/deploy/Fluently-fork/backend
    docker compose down
fi

# Crontab entry: 0 23 * * * /home/deploy/shutdown-staging.sh
# Startup: 0 8 * * * cd /home/deploy/Fluently-fork/backend && docker compose up -d
```
