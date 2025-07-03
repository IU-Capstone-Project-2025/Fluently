# Local Installation Guide (Development/Testing)

This guide helps you run Fluently on your local machine for development or grading using pre-built Docker images. **No building required** - all images are pulled from Docker Hub!

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)
- Internet connection to download pre-built images

## 🚀 Quick Start for Teaching Assistants

**Super fast setup** using pre-built images (no compilation needed!):

```bash
# 1. Clone repository
git clone https://github.com/FluentlyOrg/Fluently-fork.git
cd Fluently-fork

# 2. Copy environment configuration
cp .env.example .env

# 3. Create Docker volumes
docker volume create fluently_postgres_data
docker volume create fluently_redis_data

# 4. Check for port conflicts (optional but recommended)
echo "Checking for port conflicts..."
netstat -tulpn | grep -E ":(80|443|5432|6379|8070|8001)" || echo "No conflicts found"

# 5. Start services with pre-built images (fast!)
docker compose -f docker-compose.production.yml up -d

# 6. Wait 2-3 minutes for services to start, then access:
# - Backend API: http://localhost:8070/health  
# - Swagger UI: http://localhost:8070/swagger/
# - Frontend: http://localhost:80/
# - ML API: http://localhost:8001/health
```

> **✨ Benefits of pre-built images:**
> - **Super fast**: 2-3 minutes instead of 15+ minutes
> - **No compilation**: No Go, Python, or Node.js build process
> - **Consistent**: Same images used in production
> - **Reliable**: Pre-tested and quality-checked images

## 🔄 Alternative: Development Build (For Contributors)

If you need to modify code and build locally:

```bash
# Use the original build process
docker compose -f docker-compose.yml up -d --build
```

This takes 15+ minutes but allows code modifications.
> - We temporarily stop conflicting services - simpler than port remapping!

---


## 🌐 Access Services

### Core Application
- **Backend API:** [http://localhost:8070/api/v1/](http://localhost:8070/api/v1/)
- **Backend Health:** [http://localhost:8070/health](http://localhost:8070/health)
- **Swagger UI:** [http://localhost:8070/swagger/](http://localhost:8070/swagger/)
- **Frontend Website:** [http://localhost:80/](http://localhost:80/)
- **ML API:** [http://localhost:8001/health](http://localhost:8001/health)

### Database Access
- **PostgreSQL:** localhost:5432 (credentials in .env file)
- **Redis:** localhost:6379

## 🛠️ Management Commands

```bash
# View running services
docker compose -f docker-compose.production.yml ps

# View logs
docker compose -f docker-compose.production.yml logs -f

# Stop services
docker compose -f docker-compose.production.yml down

# Update to latest images
docker compose -f docker-compose.production.yml pull
docker compose -f docker-compose.production.yml up -d

# Clean up everything
docker compose -f docker-compose.production.yml down --volumes
docker system prune -a  # WARNING: Removes all unused Docker data
```

## 🐛 Troubleshooting

### Common Issues

**Port conflicts:**
```bash
# Check what's using the ports
sudo netstat -tulpn | grep -E ":(80|443|5432|6379|8070|8001)"

# Stop conflicting services
sudo systemctl stop postgresql apache2 nginx
```

**Services not starting:**
```bash
# Check service logs
docker compose -f docker-compose.production.yml logs [service-name]

# Restart specific service
docker compose -f docker-compose.production.yml restart [service-name]
```

**Database connection issues:**
```bash
# Check if PostgreSQL is healthy
docker compose -f docker-compose.production.yml exec postgres pg_isready

# Reset database (WARNING: Deletes all data)
docker compose -f docker-compose.production.yml down -v
docker volume rm fluently_postgres_data fluently_redis_data
# Then restart setup process
```

## 📱 Mobile Apps Integration

### Android App
- See [android-app/README.md](../android-app/README.md) for setup
- Set API base URL to `http://localhost:8070/api/v1/` in app settings

### iOS App  
- Open `ios-app/Fluently/Fluently.xcodeproj` in Xcode
- Update API base URL to `http://localhost:8070/api/v1/` for local development

## 🏗️ For Developers

If you need to modify code and test changes:

1. **Fork the repository** and create your branch
2. **Make your changes** to the code
3. **Push to your fork** - this will trigger the CI/CD pipeline
4. **Wait for images to build** in GitHub Actions
5. **Update image tags** in `docker-compose.production.yml` to use your branch
6. **Test locally** with your custom images

```bash
# Example: Using images from your fork
# Edit docker-compose.production.yml and change:
# image: docker.io/fluentlyorg/fluently-backend:latest-develop
# To:
# image: docker.io/yourusername/fluently-backend:your-branch-latest
```
