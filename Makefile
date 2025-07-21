# ===========================================
# FLUENTLY LOCAL DEVELOPMENT SETUP
# =======================================
help:                     ## Show this help message
	@echo "Fluently Local Development Setup"
	@echo "Using pre-built Docker images for fast setup"
	@echo ""
	@echo "Quick Start:"
	@echo "  1. make setup-local    # Setup environment and volumes"
	@echo "  2. make check-ports    # Check for port conflicts"
	@echo "  3. make run-local      # Start all services"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ===========================================
# SETUP COMMANDS
# ===========================================

setup-local: setup-env setup-volumes setup-thesaurus  ## Complete local setup (environment + volumes + thesaurus)
	@echo " Local development environment setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. make check-ports    # Check for port conflicts"
	@echo "  2. make run-local      # Start all services"
	@echo ""
	@echo "  IMPORTANT: Stop local PostgreSQL to avoid port conflicts:"
	@echo "   sudo systemctl stop postgresql"

setup-env:                ## Setup environment files
	@echo " Setting up environment files..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo " Created root .env from example"; \
	else \
		echo " Root .env already exists"; \
	fi

setup-volumes:            ## Create required Docker volumes
	@echo " Creating Docker volumes..."
	@docker volume create fluently_pgdata_safe || true
	@docker volume create fluently_grafana_data_external || true
	@docker volume create fluently_prometheus_data_external || true
	@echo " All volumes created"

setup-thesaurus:          ## Copy thesaurus data from deploy server
	@echo " Setting up thesaurus..."
	@if [ ! -f analysis/thesaurus/result.csv ]; then \
		echo " Copying thesaurus data (result.csv) from deploy server..."; \
		if scp deploy@45.156.22.159:/home/deploy/Fluently-fork/analysis/thesaurus/result.csv analysis/thesaurus/; then \
			echo " Thesaurus data copied successfully."; \
		else \
			echo " Failed to copy thesaurus data. Please check your SSH connection and permissions."; \
			echo "   Command: scp deploy@45.156.22.159:/home/deploy/Fluently-fork/analysis/thesaurus/result.csv analysis/thesaurus/"; \
		fi; \
	else \
		echo " Thesaurus data (result.csv) already exists."; \
	fi

# ===========================================
# MAIN COMMANDS
# ===========================================

run-local: check-env pull-images       ## Start all services with pre-built images
	@echo " Starting Fluently with pre-built images..."
	@echo "This is much faster than building locally!"
	docker compose -f docker-compose-local.yml up -d
	@echo ""
	@echo " All services started!"
	@echo ""
	@echo " Access your services:"
	@echo "  - Swagger UI:  http://localhost:8070/swagger/"
	@echo "  - Directus admin panel:     http://localhost:8055/"
	@echo "  - Distractor API:     http://localhost:8001/docs"
	@echo "  - Thesaurus API:     http://localhost:8002/docs"
	@echo "  - LLM API:     http://localhost:8003/docs"
	@echo ""
	@echo " Monitor with: make logs"

pull-images:              ## Pull latest pre-built images from Docker Hub
	@echo "� Pulling latest pre-built images..."
	@docker pull docker.io/fluentlyorg/fluently-backend:latest-develop
	@docker pull docker.io/fluentlyorg/fluently-telegram-bot:latest-develop
	@docker pull docker.io/fluentlyorg/fluently-ml-api:latest-develop
	@docker pull docker.io/fluentlyorg/fluently-nginx:latest-develop
	@echo " All images updated!"

check-ports:              ## Check for common port conflicts
	@echo " Checking for port conflicts..."
	@echo ""
	@if netstat -tlnp 2>/dev/null | grep :5432 > /dev/null; then \
		echo " Port 5432 (PostgreSQL) is in use"; \
		echo "   Fix: sudo systemctl stop postgresql"; \
		echo ""; \
	else \
		echo " Port 5432 (PostgreSQL) is free"; \
	fi
	@if netstat -tlnp 2>/dev/null | grep :80 > /dev/null; then \
		echo " Port 80 (Web server) is in use"; \
		echo "   Fix: sudo systemctl stop apache2 nginx"; \
		echo ""; \
	else \
		echo " Port 80 (Web server) is free"; \
	fi
	@if netstat -tlnp 2>/dev/null | grep :3000 > /dev/null; then \
		echo " Port 3000 (Grafana) is in use"; \
		echo "   Fix: sudo systemctl stop grafana-server"; \
		echo ""; \
	else \
		echo " Port 3000 (Grafana) is free"; \
	fi
	@if netstat -tlnp 2>/dev/null | grep :8070 > /dev/null; then \
		echo " Port 8070 (Backend API) is in use"; \
		echo ""; \
	else \
		echo " Port 8070 (Backend API) is free"; \
	fi
	@echo "If any ports are in use, stop the conflicting services before running 'make run-local'"

stop-local:               ## Stop all services
	@echo " Stopping all services..."
	docker compose down
	@echo " All services stopped"

restart:                  ## Restart all services (with fresh images)
	@echo " Restarting with latest images..."
	make stop-local
	make pull-images
	make run-local

# ===========================================
# MONITORING & DEBUGGING
# ===========================================

logs:                     ## Show logs from all services
	docker compose logs -f

logs-backend:             ## Show backend logs only
	docker compose logs -f backend

logs-ml-api:              ## Show ML API logs only
	docker compose logs -f ml-api

status:                   ## Show status of all services
	@echo " Service Status:"
	@docker compose ps

health:                   ## Check health of all services
	@echo " Health Checks:"
	@echo ""
	@echo "Backend API:"
	@curl -s http://localhost:8070/health || echo " Backend not responding"
	@echo ""
	@echo "ML API:"
	@curl -s http://localhost:8001/health || echo " ML API not responding"
	@echo ""
	@echo "Frontend:"
	@curl -s -o /dev/null -w "%%{http_code}" http://localhost/ | grep -q "200" && echo " Frontend OK" || echo " Frontend not responding"

# ===========================================
# TESTING
# ===========================================

test-backend: run-test-db ## Run backend tests with test database
	@echo " Running backend tests..."
	@cd backend && \
	export DB_HOST=localhost && \
	export DB_PORT=5433 && \
	export DB_USER=test_user && \
	export DB_NAME=test_db && \
	export DB_PASSWORD=test_password && \
	go test -v -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./...
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo " Tests complete. Coverage report: backend/coverage.html"

run-test-db:              ## Start test database
	docker compose -f docker-compose.test.yml up -d test_db

stop-test-db:             ## Stop test database
	docker compose -f docker-compose.test.yml down --volumes

# ===========================================
# CLEANUP
# ===========================================

clean:                    ## Stop services and remove volumes (DESTRUCTIVE!)
	@echo "  This will remove all data including databases!"
	@read -p "Are you sure? (y/N): " confirm && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker compose down --volumes; \
		docker volume rm fluently_pgdata_safe fluently_grafana_data_external fluently_prometheus_data_external 2>/dev/null || true; \
		echo " Cleanup complete"; \
	else \
		echo " Cleanup cancelled"; \
	fi

clean-images:             ## Remove all Fluently Docker images
	@echo " Removing all Fluently Docker images..."
	@docker images | grep fluently | awk '{print $$3}' | xargs docker rmi -f 2>/dev/null || true
	@echo " Images cleaned"

# ===========================================
# DEVELOPMENT UTILITIES
# ===========================================

generate-docs:            ## Generate API documentation
	@echo " Generating API documentation..."
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "Installing swag..."; \
		cd backend && go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@cd backend && $(HOME)/go/bin/swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal
	@echo " Documentation generated in backend/docs/"

update:                   ## Update to latest images and restart
	@echo " Updating to latest version..."
	make stop-local
	make pull-images
	make run-local
	@echo " Updated to latest version!"

# ===========================================
# HELPER FUNCTIONS
# ===========================================

check-env:
	@if [ ! -f .env ]; then \
		echo " Environment files (.env) not found!"; \
		echo "   Please run 'make setup-local' first."; \
		exit 1; \
	else \
		echo " Environment files found."; \
	fi

# ===========================================
# SYSTEM SERVICES MANAGEMENT
# ===========================================

stop-conflicting-services: ## Stop common conflicting system services
	@echo " Stopping conflicting system services..."
	@sudo systemctl stop postgresql 2>/dev/null || echo "PostgreSQL not running"
	@sudo systemctl stop apache2 2>/dev/null || echo "Apache2 not running"
	@sudo systemctl stop nginx 2>/dev/null || echo "Nginx not running"
	@sudo systemctl stop grafana-server 2>/dev/null || echo "Grafana not running"
	@echo " Conflicting services stopped"

start-system-services:    ## Restart system services after development
	@echo " Starting system services..."
	@sudo systemctl start postgresql 2>/dev/null || echo "PostgreSQL not installed"
	@sudo systemctl start apache2 2>/dev/null || echo "Apache2 not installed"
	@sudo systemctl start nginx 2>/dev/null || echo "Nginx not installed"
	@sudo systemctl start grafana-server 2>/dev/null || echo "Grafana not installed"
	@echo " System services restarted"


# Run SonarScanner (requires SONAR_TOKEN environment variable)
sonar-scan:
	@echo " Running SonarScanner analysis..."
	@if [ -z "$(SONAR_TOKEN)" ]; then \
		echo " SONAR_TOKEN environment variable is required"; \
		echo "Set it with: export SONAR_TOKEN=your_token_here"; \
		exit 1; \
	fi
	sonar-scanner -Dsonar.token=$(SONAR_TOKEN)

# Combined quality check and scan
quality-scan: code-quality sonar-scan

# ===========================================
# HELP
# ===========================================

help:
	@echo "Available commands:"
	@echo ""
	@echo "   Local Development Setup (For Teaching Assistants):"
	@echo "    setup-local         - Complete local setup (.env and volumes)"
	@echo "    setup-env           - Copy .env.example to .env"
	@echo "    setup-volumes       - Create all required local Docker volumes"
	@echo "    check-ports         - Check for common port conflicts"
	@echo ""
	@echo "   Local Development (Run after setup-local):"
	@echo "    run-local           - Start all services locally"
	@echo "    run-local-core      - Start core services locally (postgres, backend, ml-api, nginx)"
	@echo "    build-ml-api-fast   - Build ML API with optimized Docker settings"
	@echo ""
	@echo "    IMPORTANT: Stop local PostgreSQL first:"
	@echo "    sudo systemctl stop postgresql"
	@echo ""
	@echo "   Development:"
	@echo "    generate-docs       - Generate API documentation with Swagger"
	@echo "    test-backend        - Run backend tests with coverage"
	@echo "    run-backend         - Start backend with dependencies and air for hot reload"
	@echo "    run-telegram-bot    - Start telegram bot"
	@echo "    run-ml-api          - Start ML API service"
	@echo ""
	@echo "   Production (Requires external volumes):"
	@echo "    run-production      - Start all services (production mode)"
	@echo "    run-core            - Start core services (backend, db, ml-api, nginx)"
	@echo "    run-monitoring      - Start monitoring stack"
	@echo ""
	@echo "   Testing:"
	@echo "    run-test-db     - Start test database"
	@echo "    stop-test-db    - Stop test database"
	@echo ""
	@echo "   Maintenance:"
	@echo "    stop           - Stop all services"
	@echo "    clean          - Clean up volumes and containers"
	@echo "    logs SERVICE=<n> - View logs for service"
	@echo "    restart SERVICE=<n> - Restart service"
	@echo ""
	@echo "   Database:"
	@echo "    db-shell       - Access main database"
	@echo "    test-db-shell  - Access test database"
	@echo ""
	@echo "  � Code Quality:"
	@echo "    install-sonar-scanner - Install SonarScanner CLI"
	@echo "    code-quality         - Generate coverage reports"
	@echo "    sonar-scan           - Run SonarScanner (requires SONAR_TOKEN)"
	@echo "    quality-scan         - Run coverage + SonarScanner"
	@echo ""
	@echo "  � Quick Start for Local Development:"
	@echo "    make setup-local"
	@echo "    make check-ports     # Check for conflicts"
	@echo "    make build-ml-api-fast"
	@echo "    make run-local"
	@echo ""
	@echo "  Production:"
	@echo "    run-production  - Start all services"
	@echo "    run-core        - Start core services (backend, db, ml-api, nginx)"
	@echo "    run-monitoring  - Start monitoring stack"
	@echo ""
	@echo "  Testing:"
	@echo "    run-test-db     - Start test database"
	@echo "    stop-test-db    - Stop test database"
	@echo ""
	@echo "  Maintenance:"
	@echo "    stop           - Stop all services"
	@echo "    clean          - Clean up volumes and containers"
	@echo "    logs SERVICE=<name> - View logs for service"
	@echo "    restart SERVICE=<name> - Restart service"
	@echo ""
	@echo "  Database:"
	@echo "    db-shell       - Access main database"
	@echo "    test-db-shell  - Access test database"
	@echo ""
	@echo "  Fast Build:"
	@echo "    build-fast     - Build services with optimized caching"
	@echo "    build-service  - Build specific service with caching"
	@echo "    build-ml-api-fast - Build ML API with optimized settings"
	@echo "    build-go-services - Build backend and telegram bot with Go caching"