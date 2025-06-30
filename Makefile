# ===========================================
# DEVELOPMENT COMMANDS
# ===========================================

# Generate Swagger docs and run backend with supporting services
run-backend:
	cd backend && swag init --generalInfo cmd/main.go --output docs
	docker compose up -d postgres ml-api directus
	cd backend && air

# Run only the telegram bot service
run-telegram-bot:
	docker compose up --build -d telegram-bot redis

# Run the ML API service
run-ml-api:
	docker compose up --build -d ml-api

# ===========================================
# PRODUCTION COMMANDS  
# ===========================================

# Start all production services
run-production:
	docker compose up -d

# Start core services (backend, database, ML API)
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

# Clean up all volumes and orphaned containers
clean:
	docker compose down --volumes --remove-orphans
	docker compose -f docker-compose.test.yml down --volumes --remove-orphans
	rm -rf backend/tmp/* || true
	rm -rf telegram-bot/tmp/* || true

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
# HELP
# ===========================================

help:
	@echo "Available commands:"
	@echo "  Development:"
	@echo "    run-backend     - Start backend with dependencies and air for hot reload"
	@echo "    run-telegram-bot - Start telegram bot"
	@echo "    run-ml-api      - Start ML API service"
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