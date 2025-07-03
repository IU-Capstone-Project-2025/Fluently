# ===========================================
# LOCAL DEVELOPMENT SETUP
# ===========================================

# Complete local setup - creates .env and volumes
setup-local: setup-env setup-volumes
	@echo "✅ Local development environment i	@echo "    make run-local"-local"NT: Stop local PostgreSQL to avoid port conflicts:"
	@echo "   sudo systemctl stop postgresql"
	@echo ""
	@echo "Then run: make run-local"

# Copy .env.example to .env if it doesn't exist
setup-env:
	@echo "📄 Setting up environment file..."
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env; \
		echo "✅ Created backend/.env from backend/.env.example"; \
		echo "📝 Please edit backend/.env with your local settings if needed"; \
	else \
		echo "✅ backend/.env already exists"; \
	fi

# Create all required local Docker volumes
setup-volumes:
	@echo "📦 Creating local Docker volumes..."
	@docker volume create fluently_pgdata_local || true
	@docker volume create fluently_grafana_data_local || true
	@docker volume create fluently_prometheus_data_local || true
	@docker volume create fluently_sonarqube_data_local || true
	@docker volume create fluently_model_cache_local || true
	@echo "✅ All local volumes created"

# ===========================================
# LOCAL DEVELOPMENT COMMANDS
# ===========================================

# Start all services for local development
run-local:
	@echo "🚀 Starting all services for local development..."
	@if [ ! -f backend/.env ]; then \
		echo "❌ backend/.env not found!"; \
		echo "Please run 'make setup-local' first"; \
		exit 1; \
	fi
	@echo "🔍 Checking if port 5432 is free..."
	@if netstat -tlnp 2>/dev/null | grep :5432 > /dev/null; then \
		echo "❌ Port 5432 is in use (likely local PostgreSQL)"; \
		echo "Please stop local PostgreSQL first:"; \
		echo "   sudo systemctl stop postgresql"; \
		echo "Then run this command again."; \
		exit 1; \
	fi
	docker compose -f docker-compose.yml -f docker-compose.local.yml up -d

# Start core services for local development (faster)
run-local-core:
	@echo "🚀 Starting core services for local development..."
	@if [ ! -f backend/.env ]; then \
		echo "❌ backend/.env not found!"; \
		echo "Please run 'make setup-local' first"; \
		exit 1; \
	fi
	@echo "🔍 Checking if port 5432 is free..."
	@if netstat -tlnp 2>/dev/null | grep :5432 > /dev/null; then \
		echo "❌ Port 5432 is in use (likely local PostgreSQL)"; \
		echo "Please stop local PostgreSQL first:"; \
		echo "   sudo systemctl stop postgresql"; \
		echo "Then run this command again."; \
		exit 1; \
	fi
	docker compose -f docker-compose.yml -f docker-compose.local.yml up -d postgres backend ml-api nginx

# Check for common port conflicts before starting
check-ports:
	@echo "🔍 Checking for port conflicts..."
	@if netstat -tlnp 2>/dev/null | grep :5432 > /dev/null; then \
		echo "❌ Port 5432 (PostgreSQL) is in use"; \
		echo "   Run: sudo systemctl stop postgresql"; \
	else \
		echo "✅ Port 5432 (PostgreSQL) is free"; \
	fi
	@if netstat -tlnp 2>/dev/null | grep :80 > /dev/null; then \
		echo "❌ Port 80 (Web server) is in use"; \
		echo "   Check: sudo systemctl stop apache2 nginx"; \
	else \
		echo "✅ Port 80 (Web server) is free"; \
	fi
	@if netstat -tlnp 2>/dev/null | grep :3000 > /dev/null; then \
		echo "❌ Port 3000 (Grafana) is in use"; \
		echo "   Check: sudo systemctl stop grafana-server"; \
	else \
		echo "✅ Port 3000 (Grafana) is free"; \
	fi
	@if netstat -tlnp 2>/dev/null | grep :8070 > /dev/null; then \
		echo "❌ Port 8070 (Backend) is in use"; \
	else \
		echo "✅ Port 8070 (Backend) is free"; \
	fi


# Build ML API with optimized settings
build-ml-api-fast:
	@echo "🏗️ Building ML API with optimized settings..."
	docker compose build --parallel ml-api

# ===========================================
# DEVELOPMENT COMMANDS
# ===========================================

# Generate Swagger docs and run backend with supporting services
run-backend:
	cd backend && swag init --generalInfo cmd/main.go --output docs --parseDependency --parseInternal
	docker compose up -d postgres ml-api directus
	cd backend && air

# Run only the telegram bot service
run-telegram-bot:
	docker compose up --build -d telegram-bot redis

# Run the ML API service
run-ml-api:
	docker compose up --build -d ml-api

# Generate API documentation
generate-docs:
	@echo "📚 Generating API documentation..."
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "Installing swag..."; \
		cd backend && go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@cd backend && $(HOME)/go/bin/swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal
	@echo "✅ Documentation generated in backend/docs/"

# Run tests with proper test database
test-backend:
	@echo "🧪 Running backend tests..."
	@cd backend && \
	export DB_HOST=localhost && \
	export DB_PORT=5433 && \
	export DB_USER=test_user && \
	export DB_NAME=test_db && \
	export DB_PASSWORD=test_password && \
	go test -v -coverprofile=coverage.out ./...
	@cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: backend/coverage.html"

# ===========================================
# PRODUCTION COMMANDS (Use external volumes)
# ===========================================

# Start all production services (requires external volumes)
run-production:
	docker compose up -d

# Start core services (backend, database, ML API) - production
run-core:
	docker compose up -d postgres backend ml-api nginx

# Start monitoring stack
run-monitoring:
	docker compose up -d prometheus grafana loki promtail node-exporter postgres-exporter nginx-exporter cadvisor

# ===========================================
# TESTING COMMANDS
# ===========================================

# Start test database for GitHub Actions
run-test-db:
	docker compose -f docker-compose.test.yml up -d

# Stop test database
stop-test-db:
	docker compose -f docker-compose.test.yml down

# ===========================================
# MAINTENANCE COMMANDS
# ===========================================

# Stop all services
stop:
	docker compose down
	docker compose -f docker-compose.test.yml down

# Stop local development services
stop-local:
	docker compose -f docker-compose.yml -f docker-compose.local.yml down

# Clean up all volumes and orphaned containers
clean:
	docker compose down --volumes --remove-orphans
	docker compose -f docker-compose.test.yml down --volumes --remove-orphans
	rm -rf backend/tmp/* || true
	rm -rf telegram-bot/tmp/* || true

# Clean up local development environment
clean-local:
	docker compose -f docker-compose.yml -f docker-compose.local.yml down --volumes --remove-orphans
	docker volume rm fluently_pgdata fluently_grafana_data fluently_prometheus_data fluently_sonarqube_data 2>/dev/null || true
	rm -rf backend/tmp/* || true
	rm -rf analysis/distractor_api/logs/* || true

# View logs for specific service
logs:
	docker compose logs -f $(SERVICE)

# Restart specific service
restart:
	docker compose restart $(SERVICE)

# ===========================================
# DATABASE COMMANDS
# ===========================================

# Access main database
db-shell:
	docker compose exec postgres psql -U $(shell grep DB_USER backend/.env | cut -d '=' -f2) -d $(shell grep DB_NAME backend/.env | cut -d '=' -f2)

# Access test database
test-db-shell:
	docker compose -f docker-compose.test.yml exec test-db psql -U test_user -d test_fluently_db

# ===========================================
# FAST BUILD COMMANDS
# ===========================================

# Build services with optimized caching
build-fast:
	export DOCKER_BUILDKIT=1 && export COMPOSE_DOCKER_CLI_BUILD=1 && docker compose build --parallel

# Build specific service with caching
build-service:
	export DOCKER_BUILDKIT=1 && export COMPOSE_DOCKER_CLI_BUILD=1 && docker compose build $(SERVICE)

# Build ML API with optimized settings (for local development)
build-ml-api-fast:
	export DOCKER_BUILDKIT=1 && export COMPOSE_DOCKER_CLI_BUILD=1 && docker compose build ml-api

# Build backend and telegram bot with Go caching
build-go-services:
	export DOCKER_BUILDKIT=1 && export COMPOSE_DOCKER_CLI_BUILD=1 && docker compose build --parallel backend telegram-bot
=======
# CODE QUALITY COMMANDS
# ===========================================

# Install SonarScanner CLI
install-sonar-scanner:
	@echo "📦 Installing SonarScanner CLI..."
	@if ! command -v sonar-scanner >/dev/null 2>&1; then \
		echo "Installing SonarScanner CLI..."; \
		wget -q https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-5.0.1.3006-linux.zip; \
		unzip -q sonar-scanner-cli-5.0.1.3006-linux.zip; \
		sudo mv sonar-scanner-5.0.1.3006-linux /opt/sonar-scanner; \
		sudo ln -sf /opt/sonar-scanner/bin/sonar-scanner /usr/local/bin/sonar-scanner; \
		rm sonar-scanner-cli-5.0.1.3006-linux.zip; \
		echo "✅ SonarScanner CLI installed"; \
	else \
		echo "✅ SonarScanner CLI already installed"; \
	fi

# Run code quality analysis (local)
code-quality:
	@echo "🔍 Running code quality analysis..."
	@echo "📊 Generating Go coverage..."
	@cd backend && go test -coverprofile=coverage.out ./... || true
	@echo "🐍 Generating Python coverage..."
	@cd analysis/distractor_api && python3 -m pytest --cov=. --cov-report=xml || true
	@echo "✅ Coverage reports generated"

# Run SonarScanner (requires SONAR_TOKEN environment variable)
sonar-scan:
	@echo "🔍 Running SonarScanner analysis..."
	@if [ -z "$(SONAR_TOKEN)" ]; then \
		echo "❌ SONAR_TOKEN environment variable is required"; \
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
	@echo "  🏠 Local Development Setup (For Teaching Assistants):"
	@echo "    setup-local         - Complete local setup (.env and volumes)"
	@echo "    setup-env           - Copy .env.example to .env"
	@echo "    setup-volumes       - Create all required local Docker volumes"
	@echo "    check-ports         - Check for common port conflicts"
	@echo ""
	@echo "  🚀 Local Development (Run after setup-local):"
	@echo "    run-local           - Start all services locally"
	@echo "    run-local-core      - Start core services locally (postgres, backend, ml-api, nginx)"
	@echo "    build-ml-api-fast   - Build ML API with optimized Docker settings"
	@echo ""
	@echo "  ⚠️  IMPORTANT: Stop local PostgreSQL first:"
	@echo "    sudo systemctl stop postgresql"
	@echo ""
	@echo "  🔧 Development:"
	@echo "    generate-docs       - Generate API documentation with Swagger"
	@echo "    test-backend        - Run backend tests with coverage"
	@echo "    run-backend         - Start backend with dependencies and air for hot reload"
	@echo "    run-telegram-bot    - Start telegram bot"
	@echo "    run-ml-api          - Start ML API service"
	@echo ""
	@echo "  🚀 Production (Requires external volumes):"
	@echo "    run-production      - Start all services (production mode)"
	@echo "    run-core            - Start core services (backend, db, ml-api, nginx)"
	@echo "    run-monitoring      - Start monitoring stack"
	@echo ""
	@echo "  🧪 Testing:"
	@echo "    run-test-db     - Start test database"
	@echo "    stop-test-db    - Stop test database"
	@echo ""
	@echo "  🛠️ Maintenance:"
	@echo "    stop           - Stop all services"
	@echo "    clean          - Clean up volumes and containers"
	@echo "    logs SERVICE=<n> - View logs for service"
	@echo "    restart SERVICE=<n> - Restart service"
	@echo ""
	@echo "  🗄️ Database:"
	@echo "    db-shell       - Access main database"
	@echo "    test-db-shell  - Access test database"
	@echo ""
	@echo "  � Code Quality:"
	@echo "    install-sonar-scanner - Install SonarScanner CLI"
	@echo "    code-quality         - Generate coverage reports"
	@echo "    sonar-scan           - Run SonarScanner (requires SONAR_TOKEN)"
	@echo "    quality-scan         - Run coverage + SonarScanner"
	@echo ""
	@echo "  �💡 Quick Start for Local Development:"
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