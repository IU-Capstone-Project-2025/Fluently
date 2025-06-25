#!/bin/bash

# Cloudflare Origin Certificate Installer
# This script helps install Cloudflare Origin Certificates on the server

set -e

# Certificate paths
SSL_DIR="/etc/nginx/ssl"
CERT_FILE="$SSL_DIR/cloudflare-origin.pem"
KEY_FILE="$SSL_DIR/cloudflare-origin.key"

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
    echo "Usage: $0 [install|check|remove|help]"
    echo ""
    echo "Commands:"
    echo "  install   - Install Cloudflare Origin Certificates"
    echo "  check     - Check if certificates are installed and valid"
    echo "  remove    - Remove installed certificates"
    echo "  help      - Show this help message"
    echo ""
    echo "Before running 'install', you need to:"
    echo "1. Generate Origin Certificates in Cloudflare Dashboard"
    echo "2. Have the certificate and private key ready to paste"
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

# Function to install certificates
install_certificates() {
    print_info "Installing Cloudflare Origin Certificates..."
    
    create_ssl_dir
    
    echo ""
    print_warning "You need to paste your Cloudflare Origin Certificate and Private Key"
    print_warning "Make sure to include the full PEM format including BEGIN/END lines"
    echo ""
    
    # Install certificate
    print_info "Please paste your Cloudflare Origin Certificate (press Ctrl+D when done):"
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
    print_info "Please paste your Cloudflare Origin Private Key (press Ctrl+D when done):"
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
    echo -n "$cert_content" | sudo tee "$CERT_FILE" > /dev/null
    sudo chmod 644 "$CERT_FILE"
    print_status "Certificate installed: $CERT_FILE"
    
    # Write private key file
    echo -n "$key_content" | sudo tee "$KEY_FILE" > /dev/null
    sudo chmod 600 "$KEY_FILE"
    print_status "Private key installed: $KEY_FILE"
    
    # Verify installation
    if check_certificates; then
        print_status "Certificates installed successfully!"
        echo ""
        print_info "Next steps:"
        print_info "1. Set Cloudflare SSL mode to 'Full (strict)'"
        print_info "2. Use './manage-nginx.sh production --origin-certs' to enable"
        print_info "3. Test your SSL configuration"
    else
        print_error "Certificate validation failed after installation"
        return 1
    fi
}

# Function to check certificates
check_certificates() {
    print_info "Checking Cloudflare Origin Certificates..."
    
    local issues=0
    
    # Check if files exist
    if [ ! -f "$CERT_FILE" ]; then
        print_error "Certificate file not found: $CERT_FILE"
        issues=$((issues + 1))
    else
        print_status "Certificate file exists: $CERT_FILE"
        
        # Check certificate validity
        if openssl x509 -in "$CERT_FILE" -text -noout > /dev/null 2>&1; then
            print_status "Certificate is valid"
            
            # Show certificate details
            local subject=$(openssl x509 -in "$CERT_FILE" -subject -noout | sed 's/subject=//')
            local issuer=$(openssl x509 -in "$CERT_FILE" -issuer -noout | sed 's/issuer=//')
            local not_after=$(openssl x509 -in "$CERT_FILE" -dates -noout | grep notAfter | sed 's/notAfter=//')
            
            print_info "Subject: $subject"
            print_info "Issuer: $issuer"
            print_info "Expires: $not_after"
            
            # Check if it's a Cloudflare Origin certificate
            if echo "$issuer" | grep -q "Cloudflare Origin"; then
                print_status "Confirmed: Cloudflare Origin Certificate"
            else
                print_warning "This doesn't appear to be a Cloudflare Origin Certificate"
            fi
        else
            print_error "Certificate file is invalid or corrupted"
            issues=$((issues + 1))
        fi
    fi
    
    if [ ! -f "$KEY_FILE" ]; then
        print_error "Private key file not found: $KEY_FILE"
        issues=$((issues + 1))
    else
        print_status "Private key file exists: $KEY_FILE"
        
        # Check private key validity
        if openssl rsa -in "$KEY_FILE" -check -noout > /dev/null 2>&1; then
            print_status "Private key is valid"
        elif openssl pkey -in "$KEY_FILE" -check -noout > /dev/null 2>&1; then
            print_status "Private key is valid"
        else
            print_error "Private key file is invalid or corrupted"
            issues=$((issues + 1))
        fi
    fi
    
    # Check if certificate and key match
    if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
        local cert_modulus=$(openssl x509 -in "$CERT_FILE" -modulus -noout 2>/dev/null | openssl md5)
        local key_modulus=$(openssl rsa -in "$KEY_FILE" -modulus -noout 2>/dev/null | openssl md5)
        
        if [ "$cert_modulus" = "$key_modulus" ]; then
            print_status "Certificate and private key match"
        else
            print_error "Certificate and private key do not match!"
            issues=$((issues + 1))
        fi
    fi
    
    # Check file permissions
    if [ -f "$CERT_FILE" ]; then
        local cert_perms=$(stat -c "%a" "$CERT_FILE")
        if [ "$cert_perms" = "644" ]; then
            print_status "Certificate permissions are correct (644)"
        else
            print_warning "Certificate permissions: $cert_perms (should be 644)"
        fi
    fi
    
    if [ -f "$KEY_FILE" ]; then
        local key_perms=$(stat -c "%a" "$KEY_FILE")
        if [ "$key_perms" = "600" ]; then
            print_status "Private key permissions are correct (600)"
        else
            print_warning "Private key permissions: $key_perms (should be 600)"
        fi
    fi
    
    if [ $issues -eq 0 ]; then
        print_status "All certificate checks passed!"
        return 0
    else
        print_error "Found $issues issue(s) with certificates"
        return 1
    fi
}

# Function to remove certificates
remove_certificates() {
    print_warning "This will remove Cloudflare Origin Certificates"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Removing certificates..."
        
        if [ -f "$CERT_FILE" ]; then
            sudo rm "$CERT_FILE"
            print_status "Removed: $CERT_FILE"
        fi
        
        if [ -f "$KEY_FILE" ]; then
            sudo rm "$KEY_FILE"
            print_status "Removed: $KEY_FILE"
        fi
        
        # Remove SSL directory if empty
        if [ -d "$SSL_DIR" ] && [ ! "$(ls -A $SSL_DIR)" ]; then
            sudo rmdir "$SSL_DIR"
            print_status "Removed empty directory: $SSL_DIR"
        fi
        
        print_status "Certificates removed successfully"
        print_warning "Remember to switch back to Full SSL mode in nginx and Cloudflare"
    else
        print_info "Operation cancelled"
    fi
}

# Main script logic
case "${1:-}" in
    "install")
        install_certificates
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
