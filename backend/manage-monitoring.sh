#!/bin/bash

# Fluently Monitoring Stack Management Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
ZEROTIER_IP="10.243.92.227"
COMPOSE_FILE="docker-compose.yml"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    
    # Check if Docker Compose is installed
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
    
    # Check if .env file exists
    if [ ! -f .env ]; then
        log_warn ".env file not found. Creating from .env.example..."
        if [ -f .env.example ]; then
            cp .env.example .env
            log_warn "Please edit .env file with your specific values before starting services"
        else
            log_error ".env.example file not found"
            exit 1
        fi
    fi
    
    # Check ZeroTier connection
    if command -v zerotier-cli &> /dev/null; then
        if zerotier-cli info &> /dev/null; then
            log_info "ZeroTier is running"
        else
            log_warn "ZeroTier is not running or not properly configured"
        fi
    else
        log_warn "ZeroTier CLI not found. Make sure ZeroTier is installed and configured"
    fi
}

check_ports() {
    log_info "Checking if ports are available on $ZEROTIER_IP..."
    
    ports=(3000 3100 9000 9090 9100 8055 5432 8070)
    
    for port in "${ports[@]}"; do
        if netstat -tlnp 2>/dev/null | grep -q "$ZEROTIER_IP:$port"; then
            log_warn "Port $port is already in use on $ZEROTIER_IP"
        fi
    done
}

start_monitoring() {
    log_info "Starting monitoring stack..."
    
    check_prerequisites
    check_ports
    
    # Pull latest images
    log_info "Pulling latest Docker images..."
    docker-compose pull
    
    # Start services
    log_info "Starting services..."
    docker-compose up -d
    
    # Wait a moment for services to start
    sleep 10
    
    # Check service health
    check_services
    
    log_info "Monitoring stack started successfully!"
    show_urls
}

stop_monitoring() {
    log_info "Stopping monitoring stack..."
    docker-compose down
    log_info "Monitoring stack stopped"
}

restart_monitoring() {
    log_info "Restarting monitoring stack..."
    stop_monitoring
    start_monitoring
}

check_services() {
    log_info "Checking service health..."
    
    services=("prometheus" "grafana" "loki" "promtail" "sonarqube" "node-exporter")
    
    for service in "${services[@]}"; do
        if docker-compose ps | grep -q "$service.*Up"; then
            log_info "✓ $service is running"
        else
            log_error "✗ $service is not running"
        fi
    done
}

show_logs() {
    service=${1:-""}
    
    if [ -z "$service" ]; then
        log_info "Showing logs for all services..."
        docker-compose logs -f
    else
        log_info "Showing logs for $service..."
        docker-compose logs -f "$service"
    fi
}

show_urls() {
    log_info "Service URLs (accessible via ZeroTier):"
    echo "  Grafana:    http://$ZEROTIER_IP:3000 (admin/[GRAFANA_ADMIN_PASSWORD])"
    echo "  Prometheus: http://$ZEROTIER_IP:9090"
    echo "  SonarQube:  http://$ZEROTIER_IP:9000 (admin/admin - change on first login)"
    echo "  Loki:       http://$ZEROTIER_IP:3100"
    echo "  Directus:   http://$ZEROTIER_IP:8055"
    echo "  PostgreSQL: $ZEROTIER_IP:5432"
    echo ""
    echo "Note: Ensure you are connected to the ZeroTier network to access these services"
}

backup_data() {
    backup_dir="./backups"
    backup_file="monitoring-backup-$(date +%Y%m%d-%H%M%S).tar.gz"
    
    log_info "Creating backup..."
    
    mkdir -p "$backup_dir"
    
    # Stop services for consistent backup
    docker-compose down
    
    # Create backup
    sudo tar czf "$backup_dir/$backup_file" \
        /var/lib/docker/volumes/backend_prometheus_data \
        /var/lib/docker/volumes/backend_grafana_data \
        /var/lib/docker/volumes/backend_loki_data \
        /var/lib/docker/volumes/backend_sonarqube_data \
        /var/lib/docker/volumes/backend_sonarqube_extensions \
        /var/lib/docker/volumes/backend_sonarqube_logs 2>/dev/null || true
    
    # Restart services
    docker-compose up -d
    
    log_info "Backup created: $backup_dir/$backup_file"
}

show_help() {
    echo "Fluently Monitoring Stack Management"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start       Start the monitoring stack"
    echo "  stop        Stop the monitoring stack"
    echo "  restart     Restart the monitoring stack"
    echo "  status      Check service status"
    echo "  logs [SVC]  Show logs (optionally for specific service)"
    echo "  urls        Show service URLs"
    echo "  backup      Create backup of monitoring data"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start                    # Start all services"
    echo "  $0 logs grafana            # Show Grafana logs"
    echo "  $0 status                  # Check all services"
}

# Main script logic
case "${1:-help}" in
    start)
        start_monitoring
        ;;
    stop)
        stop_monitoring
        ;;
    restart)
        restart_monitoring
        ;;
    status)
        check_services
        ;;
    logs)
        show_logs "$2"
        ;;
    urls)
        show_urls
        ;;
    backup)
        backup_data
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
