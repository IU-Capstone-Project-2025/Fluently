# Cloudflare SSL Setup Guide for Fluently

This guide covers setting up SSL/TLS with Cloudflare for both production (`fluently-app.ru`) and staging (`fluently-app.online`) environments.

## Overview

Cloudflare offers different SSL modes:

1. **Flexible**: Browser ↔ HTTPS ↔ Cloudflare ↔ HTTP ↔ Server (least secure)
2. **Full**: Browser ↔ HTTPS ↔ Cloudflare ↔ HTTP/HTTPS ↔ Server (good)
3. **Full (Strict)**: Browser ↔ HTTPS ↔ Cloudflare ↔ HTTPS ↔ Server with valid cert (best)

## Current Configuration

Your Nginx configurations are set up to support both:
- **Full SSL mode**: Using `nginx.conf` and `nginx-staging.conf` (no local certificates)
- **Full (Strict) SSL mode**: Using `nginx-origin-certs.conf` and `nginx-staging-origin-certs.conf` (with Origin Certificates)

## Option 1: Full SSL Mode (Recommended for simplicity)

This is already configured and working. Cloudflare handles SSL termination, and your server runs without local certificates.

### Cloudflare Dashboard Settings:
1. Go to your domain's SSL/TLS settings
2. Set SSL/TLS encryption mode to **"Full"**
3. Enable "Always Use HTTPS"
4. Enable "Automatic HTTPS Rewrites"

### Server Configuration:
```bash
# Switch to Full SSL mode (default)
./manage-nginx.sh production
./manage-nginx.sh staging
```

## Option 2: Full (Strict) SSL Mode (Recommended for maximum security)

This provides end-to-end encryption between Cloudflare and your server.

### Step 1: Generate Cloudflare Origin Certificates

1. **In Cloudflare Dashboard:**
   - Go to SSL/TLS → Origin Server
   - Click "Create Certificate"
   - Select "Let Cloudflare generate a private key and a CSR"
   - Set hostnames: `*.fluently-app.ru, fluently-app.ru, *.fluently-app.online, fluently-app.online`
   - Choose key type: RSA (2048) or ECDSA (P-256)
   - Set certificate validity: 15 years
   - Click "Create"

2. **Copy the certificates:**
   - Copy the **Certificate** (including `-----BEGIN CERTIFICATE-----` and `-----END CERTIFICATE-----`)
   - Copy the **Private Key** (including `-----BEGIN PRIVATE KEY-----` and `-----END PRIVATE KEY-----`)

### Step 2: Install Certificates on Both Servers

#### Production Server (`fluently-app.ru`):
```bash
# SSH to production server
ssh deploy@your-production-server

# Create SSL directory
sudo mkdir -p /etc/nginx/ssl
sudo chown -R root:root /etc/nginx/ssl
sudo chmod 700 /etc/nginx/ssl

# Create certificate file
sudo tee /etc/nginx/ssl/cloudflare-origin.pem > /dev/null << 'EOF'
-----BEGIN CERTIFICATE-----
[Paste your certificate here]
-----END CERTIFICATE-----
EOF

# Create private key file
sudo tee /etc/nginx/ssl/cloudflare-origin.key > /dev/null << 'EOF'
-----BEGIN PRIVATE KEY-----
[Paste your private key here]
-----END PRIVATE KEY-----
EOF

# Set proper permissions
sudo chmod 644 /etc/nginx/ssl/cloudflare-origin.pem
sudo chmod 600 /etc/nginx/ssl/cloudflare-origin.key
```

#### Staging Server (`fluently-app.online`):
```bash
# SSH to staging server
ssh deploy-staging@your-staging-server

# Create SSL directory
sudo mkdir -p /etc/nginx/ssl
sudo chown -R root:root /etc/nginx/ssl
sudo chmod 700 /etc/nginx/ssl

# Create certificate file (same certificate works for both domains)
sudo tee /etc/nginx/ssl/cloudflare-origin.pem > /dev/null << 'EOF'
-----BEGIN CERTIFICATE-----
[Paste your certificate here - same as production]
-----END CERTIFICATE-----
EOF

# Create private key file
sudo tee /etc/nginx/ssl/cloudflare-origin.key > /dev/null << 'EOF'
-----BEGIN PRIVATE KEY-----
[Paste your private key here - same as production]
-----END PRIVATE KEY-----
EOF

# Set proper permissions
sudo chmod 644 /etc/nginx/ssl/cloudflare-origin.pem
sudo chmod 600 /etc/nginx/ssl/cloudflare-origin.key
```

### Step 3: Update Cloudflare Settings

1. **In Cloudflare Dashboard:**
   - Go to SSL/TLS → Overview
   - Set encryption mode to **"Full (strict)"**
   - Enable "Always Use HTTPS"
   - Enable "Automatic HTTPS Rewrites"

### Step 4: Switch to Origin Certificate Configuration

#### Production:
```bash
# Navigate to project directory
cd /home/deploy/Fluently-fork/backend

# Switch to production with Origin certificates
./nginx-container/manage-nginx.sh production --origin-certs

# Check status
./nginx-container/manage-nginx.sh status
```

#### Staging:
```bash
# Navigate to project directory
cd /home/deploy-staging/Fluently-fork/backend

# Switch to staging with Origin certificates
./nginx-container/manage-nginx.sh staging --origin-certs

# Check status
./nginx-container/manage-nginx.sh status
```

## Docker Compose Integration

The Origin certificates need to be mounted in the Nginx container. Add this to your `docker-compose.yml`:

```yaml
services:
  nginx:
    volumes:
      - /etc/nginx/ssl:/etc/nginx/ssl:ro  # Mount SSL certificates
      # ... other volumes
```

Or update the deployment scripts to copy certificates into the container.

## Verification

### Test SSL Configuration:
```bash
# Check certificate details
echo | openssl s_client -connect fluently-app.ru:443 -servername fluently-app.ru 2>/dev/null | openssl x509 -noout -text

# Check SSL grade
curl -s "https://api.ssllabs.com/api/v3/analyze?host=fluently-app.ru&publish=off&all=done"
```

### Check Nginx Configuration:
```bash
# Test configuration
docker compose exec nginx nginx -t

# Check active certificates
docker compose exec nginx openssl x509 -in /etc/nginx/ssl/cloudflare-origin.pem -text -noout
```

## Security Benefits of Full (Strict) Mode

1. **End-to-end encryption**: Data is encrypted from browser to Cloudflare and from Cloudflare to your server
2. **Certificate validation**: Cloudflare validates your server's certificate
3. **Protection against man-in-the-middle attacks**: Even if someone intercepts traffic between Cloudflare and your server
4. **Compliance**: Better for regulatory compliance requirements

## Troubleshooting

### Common Issues:

1. **Certificate mismatch**: Ensure the Origin certificate includes all domains (`*.fluently-app.ru`, `fluently-app.ru`, etc.)
2. **Permission errors**: Check file permissions (644 for cert, 600 for key)
3. **Path issues**: Verify certificate paths in Nginx config match actual file locations
4. **Docker volume mounts**: Ensure SSL directory is properly mounted in container

### Switching Back to Full Mode:
```bash
# If Origin certificates cause issues, switch back to Full mode
./manage-nginx.sh production    # Without --origin-certs flag
./manage-nginx.sh staging       # Without --origin-certs flag
```

## Management Commands

```bash
# Check current configuration
./manage-nginx.sh status

# Switch to production (Full SSL)
./manage-nginx.sh production

# Switch to production (Full Strict SSL)
./manage-nginx.sh production --origin-certs

# Switch to staging (Full SSL)
./manage-nginx.sh staging

# Switch to staging (Full Strict SSL)
./manage-nginx.sh staging --origin-certs

# Create backup
./manage-nginx.sh backup

# Restore from backup
./manage-nginx.sh restore
```

## Recommendations

1. **Start with Full SSL mode** to ensure everything works
2. **Upgrade to Full (Strict) SSL mode** for production environments
3. **Use the same Origin certificate** for both production and staging
4. **Set up monitoring** to alert on certificate expiration (15 years from creation)
5. **Test thoroughly** in staging before applying to production
