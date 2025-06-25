#!/bin/bash

# Nginx Configuration Manager for Fluently Project
# This script helps switch between production and staging nginx configurations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NGINX_DIR="$SCRIPT_DIR"
PROD_CONFIG="nginx.conf"
STAGING_CONFIG="nginx-staging.conf"
ACTIVE_CONFIG="nginx.conf"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [production|staging|status|backup|restore]"
    echo ""
    echo "Commands:"
    echo "  production  - Switch to production configuration (fluently-app.ru)"
    echo "  staging     - Switch to staging configuration (fluently-app.online)"
    echo "  status      - Show current configuration status"
    echo "  backup      - Create backup of current configuration"
    echo "  restore     - Restore from backup"
    echo "  help        - Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 staging                    # Switch to staging config"
    echo "  $0 production                 # Switch to production config"
    echo "  $0 status                     # Check current config"
}

# Function to backup current configuration
backup_config() {
    local backup_file="nginx.conf.backup.$(date +%Y%m%d-%H%M%S)"
    if [ -f "$NGINX_DIR/$ACTIVE_CONFIG" ]; then
        cp "$NGINX_DIR/$ACTIVE_CONFIG" "$NGINX_DIR/$backup_file"
        print_status "Configuration backed up to: $backup_file"
        return 0
    else
        print_error "No active configuration found to backup"
        return 1
    fi
}

# Function to restore from backup
restore_config() {
    local latest_backup=$(ls -t "$NGINX_DIR"/nginx.conf.backup.* 2>/dev/null | head -n1)
    if [ -n "$latest_backup" ]; then
        cp "$latest_backup" "$NGINX_DIR/$ACTIVE_CONFIG"
        print_status "Configuration restored from: $(basename $latest_backup)"
        return 0
    else
        print_error "No backup files found"
        return 1
    fi
}

# Function to switch configuration
switch_config() {
    local target_env="$1"
    local source_config=""
    
    case "$target_env" in
        "production")
            source_config="$PROD_CONFIG"
            ;;
        "staging")
            source_config="$STAGING_CONFIG"
            ;;
        *)
            print_error "Invalid environment: $target_env"
            return 1
            ;;
    esac
    
    if [ ! -f "$NGINX_DIR/$source_config" ]; then
        print_error "Source configuration not found: $source_config"
        return 1
    fi
    
    # Create backup before switching
    print_info "Creating backup before switching..."
    backup_config
    
    # Copy new configuration
    cp "$NGINX_DIR/$source_config" "$NGINX_DIR/$ACTIVE_CONFIG"
    print_status "Switched to $target_env configuration"
    
    # Test nginx configuration
    print_info "Testing nginx configuration..."
    if docker compose exec nginx nginx -t 2>/dev/null; then
        print_status "Nginx configuration test passed"
        
        # Reload nginx
        print_info "Reloading nginx..."
        if docker compose exec nginx nginx -s reload 2>/dev/null; then
            print_status "Nginx reloaded successfully"
        else
            print_warning "Could not reload nginx (container might not be running)"
        fi
    else
        print_warning "Nginx configuration test failed (container might not be running)"
        print_info "Configuration switched but nginx not reloaded"
    fi
    
    return 0
}

# Function to show current status
show_status() {
    print_info "Current nginx configuration status:"
    echo ""
    
    # Check if files exist
    if [ -f "$NGINX_DIR/$ACTIVE_CONFIG" ]; then
        print_status "Active config: $ACTIVE_CONFIG (exists)"
    else
        print_error "Active config: $ACTIVE_CONFIG (missing!)"
    fi
    
    if [ -f "$NGINX_DIR/$PROD_CONFIG" ]; then
        print_status "Production config: $PROD_CONFIG (exists)"
    else
        print_warning "Production config: $PROD_CONFIG (missing)"
    fi
    
    if [ -f "$NGINX_DIR/$STAGING_CONFIG" ]; then
        print_status "Staging config: $STAGING_CONFIG (exists)"
    else
        print_warning "Staging config: $STAGING_CONFIG (missing)"
    fi
    
    echo ""
    
    # Try to determine current environment by checking domain
    if [ -f "$NGINX_DIR/$ACTIVE_CONFIG" ]; then
        if grep -q "fluently-app\.ru" "$NGINX_DIR/$ACTIVE_CONFIG"; then
            print_status "Current environment: PRODUCTION (fluently-app.ru)"
        elif grep -q "fluently-app\.online" "$NGINX_DIR/$ACTIVE_CONFIG"; then
            print_status "Current environment: STAGING (fluently-app.online)"
        else
            print_warning "Current environment: UNKNOWN (custom configuration)"
        fi
    fi
    
    # Show backup files
    local backups=$(ls -t "$NGINX_DIR"/nginx.conf.backup.* 2>/dev/null | head -n5)
    if [ -n "$backups" ]; then
        echo ""
        print_info "Recent backups:"
        echo "$backups" | while read backup; do
            echo "  - $(basename "$backup")"
        done
    fi
}

# Main script logic
case "${1:-}" in
    "production")
        switch_config "production"
        ;;
    "staging")
        switch_config "staging"
        ;;
    "status")
        show_status
        ;;
    "backup")
        backup_config
        ;;
    "restore")
        restore_config
        ;;
    "help"|"--help"|"-h")
        show_usage
        ;;
    *)
        print_error "Invalid command: ${1:-}"
        echo ""
        show_usage
        exit 1
        ;;
esac
