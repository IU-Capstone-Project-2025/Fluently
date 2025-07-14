# Fluently - Local Installation Guide

# You can use guide from CONTRIBUTING.md

This guide helps you set up Fluently for local development using **pre-built Docker images** from Docker Hub. This approach is much faster than building images locally, especially for the ML API component.

## Prerequisites

- **Docker & Docker Compose**: [Install Docker](https://docs.docker.com/get-docker/)
- **Make**: Usually pre-installed on Linux/macOS
- **Git**: [Install Git](https://git-scm.com/downloads)
- **Minimum 8GB RAM** (16GB recommended for full monitoring stack)

## Quick Start (Recommended for TAs)

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

# 4. Start all services (uses pre-built images - much faster!)
make run-local

# 5. Access services (ready in ~2-3 minutes):
Swagger UI:  http://localhost:8070/swagger/
Directus admin panel:     http://localhost:8055/
Distractor API:     http://localhost:8001/docs
Thesaurus API:     http://localhost:8002/docs
LLM API:     http://localhost:8003/docs

# 6. When finished, stop the local build and restart your services

make stop-local           # Stop all services
sudo systemctl start postgresql  # If you stopped it
sudo systemctl start apache2     # If you stopped it
sudo systemctl start nginx       # If you stopped it
sudo systemctl start grafana-server  # If you stopped it
```

## Available Commands

```bash
# Setup & Management
make setup-local          # Initial setup (env files + volumes)
make check-ports          # Check for port conflicts
make run-local            # Start all services
make stop-local           # Stop all services
make restart              # Restart with latest images

# Monitoring & Debugging
make logs                 # Show all logs
make logs-backend         # Show backend logs only
make logs-ml-api          # Show ML API logs only
make status               # Show service status
make health               # Check health of all services

# Updates
make pull-images          # Pull latest images from Docker Hub
make update               # Update to latest version

# Testing
make test-backend         # Run backend tests
make run-test-db          # Start test database
make stop-test-db         # Stop test database

# Cleanup
make clean                # Remove all data (DESTRUCTIVE!)
make clean-images         # Remove Docker images

# Utilities
make help                 # Show all commands
make generate-docs        # Generate API documentation
```

### Core Services
- **Backend API** (Go) - Main REST API
- **ML API** (Python) - Distractor generation service  
- **Telegram Bot** (Go) - Bot service
- **Nginx** - Reverse proxy and frontend
- **PostgreSQL** - Main database
- **Redis** - Session storage

### Monitoring Stack  
- **Grafana** - Dashboards
- **Prometheus** - Metrics collection
- **Loki** - Log aggregation
- **Node Exporter** - System metrics

### Admin Tools
- **Directus** - CMS interface
- **cAdvisor** - Container metrics

### Customization
Edit environment files to customize:
- Database credentials
- API keys
- Service ports
- Domain settings
