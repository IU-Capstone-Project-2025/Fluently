#!/bin/bash

# Nginx Configuration Manager for Fluently Project
# This script helps switch between production and staging nginx configurations
# Supports Cloudflare integration with both "Full" and "Full (Strict)" SSL modes

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NGINX_DIR="$SCRIPT_DIR"
PROD_CONFIG="nginx.conf"
STAGING_CONFIG="nginx-staging.conf"
PROD_ORIGIN_CONFIG="nginx-origin-certs.conf"
STAGING_ORIGIN_CONFIG="nginx-staging-origin-certs.conf"
ACTIVE_CONFIG="nginx.conf"

# SSL certificate paths
SSL_DIR="/etc/nginx/ssl"
ORIGIN_CERT="$SSL_DIR/cloudflare-origin.pem"
ORIGIN_KEY="$SSL_DIR/cloudflare-origin.key"

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
    echo "Usage: $0 [production|staging|status|backup|restore] [--origin-certs]"
    echo ""
    echo "Commands:"
    echo "  production  - Switch to production configuration (fluently-app.ru)"
    echo "  staging     - Switch to staging configuration (fluently-app.online)"
    echo "  status      - Show current configuration status"
    echo "  backup      - Create backup of current configuration"
    echo "  restore     - Restore from backup"
    echo "  help        - Show this help message"
    echo ""
    echo "Options:"
    echo "  --origin-certs    Use Cloudflare Origin Certificates (Full Strict SSL)"
    echo "                    Default: Use Cloudflare Full SSL (no local certs)"
    echo ""
    echo "Cloudflare SSL Modes:"
    echo "  Full SSL:         Cloudflare <-HTTPS-> Browser, Cloudflare <-HTTP-> Server"
    echo "  Full (Strict):    Cloudflare <-HTTPS-> Browser, Cloudflare <-HTTPS-> Server"
    echo ""
    echo "Examples:"
    echo "  $0 staging                    # Switch to staging (Full SSL)"
    echo "  $0 production --origin-certs  # Switch to production (Full Strict SSL)"
    echo "  $0 status                     # Check current config and SSL mode"
}

# Function to check if Cloudflare Origin certificates exist
check_origin_certs() {
    if [ -f "$ORIGIN_CERT" ] && [ -f "$ORIGIN_KEY" ]; then
        return 0
    else
        return 1
    fi
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
    local use_origin_certs="$2"
    local source_config=""
    
    case "$target_env" in
        "production")
            if [ "$use_origin_certs" = "true" ]; then
                source_config="$PROD_ORIGIN_CONFIG"
            else
                source_config="$PROD_CONFIG"
            fi
            ;;
        "staging")
            if [ "$use_origin_certs" = "true" ]; then
                source_config="$STAGING_ORIGIN_CONFIG"
            else
                source_config="$STAGING_CONFIG"
            fi
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
    
    # Check for Origin certificates if needed
    if [ "$use_origin_certs" = "true" ]; then
        if ! check_origin_certs; then
            print_error "Cloudflare Origin Certificates not found!"
            print_warning "Expected files:"
            print_warning "  - $ORIGIN_CERT"
            print_warning "  - $ORIGIN_KEY"
            print_info "Please install Origin Certificates or use without --origin-certs flag"
            return 1
        else
            print_status "Origin certificates found, using Full (Strict) SSL mode"
        fi
    else
        print_status "Using Cloudflare Full SSL mode (no local certificates)"
    fi
    
    # Create backup before switching
    print_info "Creating backup before switching..."
    backup_config
    
    # Copy new configuration
    cp "$NGINX_DIR/$source_config" "$NGINX_DIR/$ACTIVE_CONFIG"
    
    if [ "$use_origin_certs" = "true" ]; then
        print_status "Switched to $target_env configuration with Origin Certificates"
    else
        print_status "Switched to $target_env configuration (Cloudflare Full SSL)"
    fi
    
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
    
    if [ -f "$NGINX_DIR/$PROD_ORIGIN_CONFIG" ]; then
        print_status "Production Origin config: $PROD_ORIGIN_CONFIG (exists)"
    else
        print_warning "Production Origin config: $PROD_ORIGIN_CONFIG (missing)"
    fi
    
    if [ -f "$NGINX_DIR/$STAGING_ORIGIN_CONFIG" ]; then
        print_status "Staging Origin config: $STAGING_ORIGIN_CONFIG (exists)"
    else
        print_warning "Staging Origin config: $STAGING_ORIGIN_CONFIG (missing)"
    fi
    
    echo ""
    
    # Check Cloudflare Origin Certificates
    if check_origin_certs; then
        print_status "Cloudflare Origin Certificates: Available"
    else
        print_warning "Cloudflare Origin Certificates: Not found"
        print_info "  Expected at: $ORIGIN_CERT and $ORIGIN_KEY"
    fi
    
    echo ""
    
    # Try to determine current environment and SSL mode
    if [ -f "$NGINX_DIR/$ACTIVE_CONFIG" ]; then
        if grep -q "fluently-app\.ru" "$NGINX_DIR/$ACTIVE_CONFIG"; then
            if grep -q "ssl_certificate.*cloudflare-origin" "$NGINX_DIR/$ACTIVE_CONFIG"; then
                print_status "Current environment: PRODUCTION (Full Strict SSL with Origin Certs)"
            else
                print_status "Current environment: PRODUCTION (Full SSL via Cloudflare)"
            fi
        elif grep -q "fluently-app\.online" "$NGINX_DIR/$ACTIVE_CONFIG"; then
            if grep -q "ssl_certificate.*cloudflare-origin" "$NGINX_DIR/$ACTIVE_CONFIG"; then
                print_status "Current environment: STAGING (Full Strict SSL with Origin Certs)"
            else
                print_status "Current environment: STAGING (Full SSL via Cloudflare)"
            fi
        else
            print_warning "Current environment: UNKNOWN (custom configuration)"
        fi
        
        # Check SSL configuration details
        if grep -q "listen 443 ssl" "$NGINX_DIR/$ACTIVE_CONFIG"; then
            print_info "SSL Mode: HTTPS with certificates (Full Strict)"
        elif grep -q "listen 443" "$NGINX_DIR/$ACTIVE_CONFIG"; then
            print_info "SSL Mode: HTTPS without local certificates (Full)"
        else
            print_warning "SSL Mode: HTTP only (not recommended for production)"
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
use_origin_certs="false"

# Parse command line arguments
case "${1:-}" in
    "production"|"staging")
        command="$1"
        if [ "${2:-}" = "--origin-certs" ]; then
            use_origin_certs="true"
        fi
        switch_config "$command" "$use_origin_certs"
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
