# Full Production Installation Guide

This guide explains how to deploy Fluently in a production-like environment with your own domain, SSL certificates, and all services running as in production.

---

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)
- Your own domain name (e.g., `your-domain.com`)
- (Optional) VPS or cloud server with public IP

---

## 1. Clone the Repository
```bash
git clone https://github.com/IU-Capstone-Project-2025/Fluently.git
cd Fluently
```

---

## 2. Create Environment File
Create a `.env` file in the `backend` directory:
```env
# JWT
JWT_SECRET=supersecretjwtkey
JWT_EXPIRATION=24h
REFRESH_EXPIRATION=720h

# App
APP_NAME=fluently
# Port for the backend API, recommended to use something other than 8080 to avoid conflicts
APP_PORT=8079 

# Database
DB_USER=postgres
DB_PASSWORD=securepassword
DB_HOST=postgres
DB_PORT=5432
DB_NAME=postgres

# Directus configuration
# Directus port recommended to use something other than 8080 to avoid conflicts
DIRECTUS_PORT=8078
DIRECTUS_ADMIN_EMAIL=admin@example.com
DIRECTUS_ADMIN_PASSWORD=admin
DIRECTUS_SECRET_KEY=supersecretkey

# Google OAuth
IOS_GOOGLE_CLIENT_ID=some-ios-client-id.apps.googleusercontent.com
ANDROID_GOOGLE_CLIENT_ID=some-android-client-id.apps.googleusercontent.com
WEB_GOOGLE_CLIENT_ID=some-web-client-id.apps.googleusercontent.com

```

---

## 3. Configure Nginx for Your Domain
Edit `backend/nginx-container/nginx.conf`:
```nginx
server_name your-domain.com www.your-domain.com;
```
Replace all instances of the default domain with your own.

---

## 4. Set Up SSL Certificates
- Use [Let's Encrypt](https://letsencrypt.org/) and [Certbot](https://certbot.eff.org/) to generate SSL certificates for your domain. (Is is 100% free and automatically renews each 90 days)
- Command to generate certificates:
```bash
sudo certbot certonly --nginx -d your-domain.com -d www.your-domain.com
```
- Follow the prompts to complete the certificate generation.
- Place the certificates in the correct paths as referenced in `nginx.conf`:
  - `/etc/letsencrypt/live/your-domain.com/fullchain.pem`
  - `/etc/letsencrypt/live/your-domain.com/privkey.pem`

---

## 5. Start All Services
From Fluently/backend:
```bash
docker compose up -d --build
```

---

## 6. Access Services
- **Backend API:** `https://your-domain.com/api/v1/`
- **Swagger UI:** `https://your-domain.com/swagger/`
- **Frontend (static):** `https://your-domain.com/`
- **Directus (if enabled):** `https://your-domain.com/admin/`

---

## 7. CI/CD (Optional)
For automated deployments, configure GitHub Actions with these secrets:
```yaml
DEPLOY_HOST: SSH server IP
DEPLOY_USERNAME: SSH username
DEPLOY_SSH_KEY: Private SSH key
```
See [`.github/workflows/deploy.yml`](../.github/workflows/deploy.yml) for an example workflow.

---

## 8. Stopping Services
From Fluently/backend:
```bash
docker compose down
```

---

## Notes
- This setup is intended for production or advanced testing.
- Requires a real domain and public server.
- For local testing, see [Local Installation Guide](Install_Local.md).
