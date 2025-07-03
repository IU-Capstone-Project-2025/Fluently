# Local Installation Guide (Development/Testing)

This guide helps you run Fluently on your local machine for development or grading. No domain or SSL required. All services run on `localhost` using Docker Compose.

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)
- [Make](https://www.gnu.org/software/make/) (Linux/macOS) or [Make for Windows](https://gnuwin32.sourceforge.net/packadocker compose -f docker-compose.yml -f docker-compose.local.yml up -dges/make.htm)

## 🚀 Quick Start for Teaching Assistants

Simple step-by-step process:

```bash
# 1. Clone and setup
git clone https://github.com/FluentlyOrg/Fluently-fork.git
cd Fluently-fork
make setup-local

# 2. Check for port conflicts
make check-ports

# 3. Stop conflicting services (if needed from step 2)
sudo systemctl stop postgresql  # Always needed
sudo systemctl stop apache2     # If you have Apache
sudo systemctl stop nginx       # If you have Nginx  
sudo systemctl stop grafana-server  # If you have Grafana

# 4. Build ML API with optimizations (reduces build time)
make build-ml-api-fast

# 5. Start all services
make run-local

# 6. Wait 8-10 minutes for optimized build, then access services:
# - Backend API: http://localhost:8070/health
# - Swagger UI: http://localhost:8070/swagger/
# - Frontend: http://localhost:80/
# - Database: localhost:5432 (standard port)
# - Grafana: http://localhost:3000/

# 7. When finished, restart your services
sudo systemctl start postgresql
sudo systemctl start apache2     # If you stopped it
sudo systemctl start nginx       # If you stopped it
sudo systemctl start grafana-server  # If you stopped it
```

> **Note**: 
> - The optimized ML API build takes 8-10 minutes instead of 15+ minutes
> - We temporarily stop conflicting services - simpler than port remapping!

---


##  Access Services
- **Backend API:** [http://localhost:8070/api/v1/](http://localhost:8070/api/v1/)
- **Swagger UI:** [http://localhost:8070/swagger/](http://localhost:8070/swagger/)
- **Frontend Website:** [http://localhost:80/](http://localhost:80/)
- **Database:** localhost:5432 (standard PostgreSQL port)
- **Directus CMS:** [http://localhost:8055/admin](http://localhost:8055/)
- **ML API (Internal):** [http://localhost:8001/health](http://localhost:8001/health)

### Monitoring Services (Optional)
- **Grafana:** [http://localhost:3000](http://localhost:3000) (admin/admin123)
- **Prometheus:** [http://localhost:9090](http://localhost:9090)

### Code Quality Analysis
- **SonarScanner CLI:** Use `make install-sonar-scanner` and `make quality-scan`
- **Local Coverage:** Generated in `backend/coverage.html` and `analysis/distractor_api/htmlcov/`

---


## Mobile Apps (Optional)
- **Android:** See [android-app/README.md](../android-app/README.md) for setup instructions
- **iOS:** Open `ios-app/Fluently/Fluently.xcodeproj` in Xcode
  - Set API base URL to `http://localhost:8070/` for local development
  - Or use staging: `https://fluently-app.online/api/v1/`
