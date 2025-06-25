#!/bin/bash

# Cloudflare Origin Certificate Manager for Multiple Domains
# This script helps install domain-specific Cloudflare Origin Certificates

set -e

# Certificate paths
SSL_DIR="/etc/nginx/ssl"
PROD_CERT_FILE="$SSL_DIR/fluently-app-ru.pem"
PROD_KEY_FILE="$SSL_DIR/fluently-app-ru.key"
STAGING_CERT_FILE="$SSL_DIR/fluently-app-online.pem"
STAGING_KEY_FILE="$SSL_DIR/fluently-app-online.key"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

show_usage() {
    echo "Usage: $0 [install-production|install-staging|check|remove|help]"
    echo ""
    echo "Commands:"
    echo "  install-production  - Install Origin Certificate for fluently-app.ru"
    echo "  install-staging     - Install Origin Certificate for fluently-app.online"
    echo "  check              - Check if certificates are installed and valid"
    echo "  remove             - Remove all installed certificates"
    echo "  help               - Show this help message"
    echo ""
    echo "Domain-specific certificates:"
    echo "  Production:  fluently-app.ru, *.fluently-app.ru"
    echo "  Staging:     fluently-app.online, *.fluently-app.online"
    echo ""
    echo "Before running install commands, generate Origin Certificates in Cloudflare:"
    echo "1. Go to Cloudflare Dashboard for the specific domain"
    echo "2. SSL/TLS → Origin Server → Create Certificate"
    echo "3. Set hostnames for the domain (e.g., *.fluently-app.ru, fluently-app.ru)"
}

# Function to create SSL directory
create_ssl_dir() {
    if [ ! -d "$SSL_DIR" ]; then
        print_info "Creating SSL directory: $SSL_DIR"
        sudo mkdir -p "$SSL_DIR"
        sudo chown root:root "$SSL_DIR"
        sudo chmod 700 "$SSL_DIR"
        print_status "SSL directory created"
    else
        print_status "SSL directory already exists"
    fi
}

# Function to install certificate for a specific domain
install_certificate() {
    local domain="$1"
    local cert_file="$2"
    local key_file="$3"
    
    print_info "Installing Cloudflare Origin Certificate for $domain..."
    
    create_ssl_dir
    
    echo ""
    print_warning "You need to paste your Cloudflare Origin Certificate for $domain"
    print_warning "Make sure to include the full PEM format including BEGIN/END lines"
    echo ""
    
    # Install certificate
    print_info "Please paste your Cloudflare Origin Certificate for $domain (press Ctrl+D when done):"
    echo "Expected hostnames: *.$domain, $domain"
    echo ""
    echo "Expected format:"
    echo "-----BEGIN CERTIFICATE-----"
    echo "... certificate content ..."
    echo "-----END CERTIFICATE-----"
    echo ""
    echo "Paste certificate:"
    
    local cert_content=""
    while IFS= read -r line; do
        cert_content+="$line"$'\n'
    done
    
    if [ -z "$cert_content" ]; then
        print_error "No certificate content provided"
        return 1
    fi
    
    # Validate certificate format
    if [[ ! "$cert_content" =~ -----BEGIN\ CERTIFICATE----- ]] || [[ ! "$cert_content" =~ -----END\ CERTIFICATE----- ]]; then
        print_error "Invalid certificate format. Make sure to include BEGIN and END lines."
        return 1
    fi
    
    echo ""
    print_info "Please paste your Cloudflare Origin Private Key for $domain (press Ctrl+D when done):"
    echo "Expected format:"
    echo "-----BEGIN PRIVATE KEY-----"
    echo "... private key content ..."
    echo "-----END PRIVATE KEY-----"
    echo ""
    echo "Paste private key:"
    
    local key_content=""
    while IFS= read -r line; do
        key_content+="$line"$'\n'
    done
    
    if [ -z "$key_content" ]; then
        print_error "No private key content provided"
        return 1
    fi
    
    # Validate private key format
    if [[ ! "$key_content" =~ -----BEGIN\ PRIVATE\ KEY----- ]] && [[ ! "$key_content" =~ -----BEGIN\ RSA\ PRIVATE\ KEY----- ]]; then
        print_error "Invalid private key format. Make sure to include BEGIN and END lines."
        return 1
    fi
    
    # Write certificate file
    echo -n "$cert_content" | sudo tee "$cert_file" > /dev/null
    sudo chmod 644 "$cert_file"
    print_status "Certificate installed: $cert_file"
    
    # Write private key file
    echo -n "$key_content" | sudo tee "$key_file" > /dev/null
    sudo chmod 600 "$key_file"
    print_status "Private key installed: $key_file"
    
    # Verify installation
    if check_certificate_pair "$cert_file" "$key_file" "$domain"; then
        print_status "Certificate for $domain installed successfully!"
    else
        print_error "Certificate validation failed after installation"
        return 1
    fi
}

# Function to check a certificate pair
check_certificate_pair() {
    local cert_file="$1"
    local key_file="$2"
    local domain="$3"
    
    local issues=0
    
    # Check if files exist
    if [ ! -f "$cert_file" ]; then
        print_error "Certificate file not found: $cert_file"
        return 1
    fi
    
    if [ ! -f "$key_file" ]; then
        print_error "Private key file not found: $key_file"
        return 1
    fi
    
    # Check certificate validity
    if openssl x509 -in "$cert_file" -text -noout > /dev/null 2>&1; then
        print_status "Certificate is valid: $cert_file"
        
        # Show certificate details
        local subject=$(openssl x509 -in "$cert_file" -subject -noout | sed 's/subject=//')
        local issuer=$(openssl x509 -in "$cert_file" -issuer -noout | sed 's/issuer=//')
        local not_after=$(openssl x509 -in "$cert_file" -dates -noout | grep notAfter | sed 's/notAfter=//')
        
        print_info "Subject: $subject"
        print_info "Issuer: $issuer"
        print_info "Expires: $not_after"
        
        # Check if it's a Cloudflare Origin certificate
        if echo "$issuer" | grep -q "CloudFlare Origin"; then
            print_status "Confirmed: Cloudflare Origin Certificate"
        else
            print_warning "This doesn't appear to be a Cloudflare Origin Certificate"
        fi
        
        # Check domain coverage
        local sans=$(openssl x509 -in "$cert_file" -text -noout | grep -A 1 "Subject Alternative Name" | tail -1)
        if echo "$sans" | grep -q "$domain"; then
            print_status "Certificate covers domain: $domain"
        else
            print_warning "Certificate may not cover domain: $domain"
            print_info "SANs: $sans"
        fi
    else
        print_error "Certificate file is invalid or corrupted"
        issues=$((issues + 1))
    fi
    
    # Check private key validity
    if openssl rsa -in "$key_file" -check -noout > /dev/null 2>&1; then
        print_status "Private key is valid: $key_file"
    elif openssl pkey -in "$key_file" -check -noout > /dev/null 2>&1; then
        print_status "Private key is valid: $key_file"
    else
        print_error "Private key file is invalid or corrupted"
        issues=$((issues + 1))
    fi
    
    # Check if certificate and key match
    local cert_modulus=$(openssl x509 -in "$cert_file" -modulus -noout 2>/dev/null | openssl md5)
    local key_modulus=$(openssl rsa -in "$key_file" -modulus -noout 2>/dev/null | openssl md5)
    
    if [ "$cert_modulus" = "$key_modulus" ]; then
        print_status "Certificate and private key match"
    else
        print_error "Certificate and private key do not match!"
        issues=$((issues + 1))
    fi
    
    return $issues
}

# Function to check all certificates
check_certificates() {
    print_info "Checking Cloudflare Origin Certificates..."
    echo ""
    
    local total_issues=0
    
    # Check production certificates
    print_info "=== Production Certificates (fluently-app.ru) ==="
    if check_certificate_pair "$PROD_CERT_FILE" "$PROD_KEY_FILE" "fluently-app.ru"; then
        print_status "Production certificates: OK"
    else
        print_warning "Production certificates: Issues found"
        total_issues=$((total_issues + 1))
    fi
    
    echo ""
    
    # Check staging certificates
    print_info "=== Staging Certificates (fluently-app.online) ==="
    if check_certificate_pair "$STAGING_CERT_FILE" "$STAGING_KEY_FILE" "fluently-app.online"; then
        print_status "Staging certificates: OK"
    else
        print_warning "Staging certificates: Issues found"
        total_issues=$((total_issues + 1))
    fi
    
    echo ""
    
    if [ $total_issues -eq 0 ]; then
        print_status "All certificate checks passed!"
        echo ""
        print_info "Next steps:"
        print_info "1. Set Cloudflare SSL mode to 'Full (strict)' for both domains"
        print_info "2. Use './manage-nginx.sh production --origin-certs' for production"
        print_info "3. Use './manage-nginx.sh staging --origin-certs' for staging"
        return 0
    else
        print_error "Found issues with $total_issues certificate pair(s)"
        return 1
    fi
}

# Function to remove certificates
remove_certificates() {
    print_warning "This will remove ALL Cloudflare Origin Certificates"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Removing certificates..."
        
        for file in "$PROD_CERT_FILE" "$PROD_KEY_FILE" "$STAGING_CERT_FILE" "$STAGING_KEY_FILE"; do
            if [ -f "$file" ]; then
                sudo rm "$file"
                print_status "Removed: $file"
            fi
        done
        
        # Remove SSL directory if empty
        if [ -d "$SSL_DIR" ] && [ ! "$(ls -A $SSL_DIR)" ]; then
            sudo rmdir "$SSL_DIR"
            print_status "Removed empty directory: $SSL_DIR"
        fi
        
        print_status "All certificates removed successfully"
        print_warning "Remember to switch back to Full SSL mode in nginx and Cloudflare"
    else
        print_info "Operation cancelled"
    fi
}

# Main script logic
case "${1:-}" in
    "install-production")
        install_certificate "fluently-app.ru" "$PROD_CERT_FILE" "$PROD_KEY_FILE"
        ;;
    "install-staging")
        install_certificate "fluently-app.online" "$STAGING_CERT_FILE" "$STAGING_KEY_FILE"
        ;;
    "check")
        check_certificates
        ;;
    "remove")
        remove_certificates
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
