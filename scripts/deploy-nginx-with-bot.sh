#!/bin/bash

# Comprehensive deployment script for Fluently backend + Telegram bot
# This script handles NGINX configuration with bot integration

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Fluently Full Stack Deployment (Backend + Bot)${NC}"
echo -e "${BLUE}====================================================${NC}"

# Function to check if running as root (needed for NGINX operations)
check_root() {
    if [ "$EUID" -ne 0 ]; then
        echo -e "${RED}❌ This script must be run as root for NGINX operations${NC}"
        echo -e "${YELLOW}💡 Run with: sudo $0${NC}"
        exit 1
    fi
}

# Function to load environment variables
load_env() {
    if [ -f telegram-bot/.env ]; then
        echo -e "${GREEN}✅ Loading bot environment variables${NC}"
        # Export variables from bot .env
        set -a
        source telegram-bot/.env
        set +a
    else
        echo -e "${RED}❌ Bot .env file not found${NC}"
        exit 1
    fi
    
    # Set default values if not provided
    export DOMAIN=${DOMAIN:-fluently-app.ru}
    export CERT_NAME=${CERT_NAME:-fluently-app.ru}
    export WEBHOOK_SECRET=${WEBHOOK_SECRET:-1234567890}
    
    echo -e "${YELLOW}📋 Configuration loaded:${NC}"
    echo -e "🌐 Domain: $DOMAIN"
    echo -e "🔑 Webhook Secret: $WEBHOOK_SECRET"
    echo -e "📜 Certificate: $CERT_NAME"
}

# Function to generate NGINX configuration from template
generate_nginx_config() {
    echo -e "\n${YELLOW}📝 Generating NGINX configuration...${NC}"
    
    local template_file="backend/nginx-container/nginx.conf.template"
    local output_file="/etc/nginx/sites-available/fluently-full"
    
    if [ ! -f "$template_file" ]; then
        echo -e "${RED}❌ NGINX template not found: $template_file${NC}"
        exit 1
    fi
    
    # Use envsubst to substitute environment variables
    envsubst '${DOMAIN} ${CERT_NAME} ${WEBHOOK_SECRET}' < "$template_file" > "$output_file"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ NGINX configuration generated: $output_file${NC}"
    else
        echo -e "${RED}❌ Failed to generate NGINX configuration${NC}"
        exit 1
    fi
}

# Function to setup NGINX
setup_nginx() {
    echo -e "\n${YELLOW}🔧 Setting up NGINX...${NC}"
    
    # Disable default site if exists
    if [ -L /etc/nginx/sites-enabled/default ]; then
        echo -e "${YELLOW}📴 Disabling default NGINX site${NC}"
        rm /etc/nginx/sites-enabled/default
    fi
    
    # Enable our site
    if [ ! -L /etc/nginx/sites-enabled/fluently-full ]; then
        echo -e "${YELLOW}🔗 Enabling Fluently site${NC}"
        ln -s /etc/nginx/sites-available/fluently-full /etc/nginx/sites-enabled/
    fi
    
    # Test NGINX configuration
    echo -e "${YELLOW}🧪 Testing NGINX configuration...${NC}"
    nginx -t
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ NGINX configuration is valid${NC}"
    else
        echo -e "${RED}❌ NGINX configuration has errors${NC}"
        exit 1
    fi
    
    # Reload NGINX
    echo -e "${YELLOW}🔄 Reloading NGINX...${NC}"
    systemctl reload nginx
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ NGINX reloaded successfully${NC}"
    else
        echo -e "${RED}❌ Failed to reload NGINX${NC}"
        exit 1
    fi
}

# Function to check prerequisites
check_prerequisites() {
    echo -e "\n${YELLOW}🔍 Checking prerequisites...${NC}"
    
    # Check if NGINX is installed
    if ! command -v nginx &> /dev/null; then
        echo -e "${RED}❌ NGINX is not installed${NC}"
        echo -e "${YELLOW}💡 Install with: apt update && apt install nginx${NC}"
        exit 1
    fi
    
    # Check if envsubst is installed
    if ! command -v envsubst &> /dev/null; then
        echo -e "${RED}❌ envsubst is not installed${NC}"
        echo -e "${YELLOW}💡 Install with: apt install gettext-base${NC}"
        exit 1
    fi
    
    # Check if SSL certificates exist
    if [ ! -f "/etc/nginx/ssl/${CERT_NAME}.pem" ] || [ ! -f "/etc/nginx/ssl/${CERT_NAME}.key" ]; then
        echo -e "${YELLOW}⚠️  SSL certificates not found at /etc/nginx/ssl/${CERT_NAME}.*${NC}"
        echo -e "${YELLOW}💡 Make sure your Cloudflare Origin Certificates are installed${NC}"
        echo -e "${YELLOW}💡 Or update CERT_NAME in .env if using different certificates${NC}"
        
        # Ask user if they want to continue anyway
        read -p "Continue anyway? (y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    else
        echo -e "${GREEN}✅ SSL certificates found${NC}"
    fi
    
    # Check if Redis is running
    if ! systemctl is-active --quiet redis; then
        echo -e "${YELLOW}⚠️  Redis is not running${NC}"
        echo -e "${YELLOW}💡 Starting Redis...${NC}"
        systemctl start redis
    fi
    
    echo -e "${GREEN}✅ Prerequisites check completed${NC}"
}

# Function to deploy bot
deploy_bot() {
    echo -e "\n${YELLOW}🤖 Deploying Telegram bot...${NC}"
    
    # Change to bot directory and run deployment
    cd telegram-bot
    
    # Make sure the deployment script is executable
    chmod +x scripts/deploy-bot.sh
    
    # Run bot deployment (skip webhook registration since we'll do it separately)
    sudo -u $(logname) ./scripts/deploy-bot.sh --skip-webhook
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Bot deployment completed${NC}"
    else
        echo -e "${RED}❌ Bot deployment failed${NC}"
        exit 1
    fi
    
    cd ..
}

# Function to register webhook
register_webhook() {
    echo -e "\n${YELLOW}📡 Registering Telegram webhook...${NC}"
    
    cd telegram-bot
    
    # Make sure the webhook script is executable
    chmod +x scripts/register-webhook.sh
    
    # Run webhook registration as regular user
    sudo -u $(logname) ./scripts/register-webhook.sh
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Webhook registration completed${NC}"
    else
        echo -e "${RED}❌ Webhook registration failed${NC}"
        exit 1
    fi
    
    cd ..
}

# Function to test deployment
test_deployment() {
    echo -e "\n${YELLOW}🧪 Testing deployment...${NC}"
    
    # Test backend health (if backend is running)
    echo -e "${YELLOW}🔍 Testing backend health...${NC}"
    if curl -s "https://${DOMAIN}/health" | grep -q "ok\|healthy"; then
        echo -e "${GREEN}✅ Backend health endpoint responding${NC}"
    else
        echo -e "${YELLOW}⚠️  Backend health endpoint not responding (backend may not be running)${NC}"
    fi
    
    # Test bot health
    echo -e "${YELLOW}🔍 Testing bot health...${NC}"
    if curl -s "https://${DOMAIN}/bot/health" | grep -q "healthy"; then
        echo -e "${GREEN}✅ Bot health endpoint responding${NC}"
    else
        echo -e "${RED}❌ Bot health endpoint not responding${NC}"
    fi
    
    # Test webhook endpoint (should return 405 Method Not Allowed for GET)
    echo -e "${YELLOW}🔍 Testing webhook endpoint...${NC}"
    webhook_response=$(curl -s -o /dev/null -w "%{http_code}" "https://${DOMAIN}/webhook")
    if [ "$webhook_response" = "405" ]; then
        echo -e "${GREEN}✅ Webhook endpoint responding (405 Method Not Allowed is expected for GET)${NC}"
    else
        echo -e "${RED}❌ Webhook endpoint not responding correctly (got $webhook_response)${NC}"
    fi
}

# Function to show status
show_status() {
    echo -e "\n${BLUE}📊 Deployment Status${NC}"
    echo -e "${BLUE}===================${NC}"
    
    # NGINX status
    if systemctl is-active --quiet nginx; then
        echo -e "${GREEN}✅ NGINX: Running${NC}"
    else
        echo -e "${RED}❌ NGINX: Not running${NC}"
    fi
    
    # Redis status
    if systemctl is-active --quiet redis; then
        echo -e "${GREEN}✅ Redis: Running${NC}"
    else
        echo -e "${RED}❌ Redis: Not running${NC}"
    fi
    
    # Bot status
    if pgrep -f "telegram-bot" > /dev/null; then
        echo -e "${GREEN}✅ Telegram Bot: Running${NC}"
    else
        echo -e "${RED}❌ Telegram Bot: Not running${NC}"
    fi
    
    echo -e "\n${BLUE}🔗 Useful URLs:${NC}"
    echo -e "🌐 Main site: https://${DOMAIN}/"
    echo -e "🔍 Backend health: https://${DOMAIN}/health"
    echo -e "🤖 Bot health: https://${DOMAIN}/bot/health"
    echo -e "📡 Webhook: https://${DOMAIN}/webhook"
    
    echo -e "\n${BLUE}📋 Management Commands:${NC}"
    echo -e "• Check bot logs: tail -f telegram-bot/logs/bot.log"
    echo -e "• Check NGINX logs: tail -f /var/log/nginx/access.log"
    echo -e "• Restart bot: pkill -f telegram-bot && cd telegram-bot && ./bin/telegram-bot &"
    echo -e "• Reload NGINX: systemctl reload nginx"
}

# Main deployment function
main() {
    echo -e "\n${BLUE}🔄 Starting full deployment...${NC}"
    
    # Step 1: Check prerequisites
    check_prerequisites
    
    # Step 2: Load environment
    load_env
    
    # Step 3: Generate NGINX configuration
    generate_nginx_config
    
    # Step 4: Setup NGINX
    setup_nginx
    
    # Step 5: Deploy bot
    deploy_bot
    
    # Step 6: Register webhook
    register_webhook
    
    # Step 7: Test deployment
    test_deployment
    
    # Step 8: Show status
    show_status
    
    echo -e "\n${GREEN}🎉 Full deployment completed successfully!${NC}"
    echo -e "${YELLOW}💡 Send a message to your bot to test end-to-end functionality${NC}"
}

# Handle command line arguments
case "${1}" in
    --nginx-only)
        check_root
        load_env
        generate_nginx_config
        setup_nginx
        ;;
    --bot-only)
        load_env
        deploy_bot
        ;;
    --webhook-only)
        load_env
        register_webhook
        ;;
    --test-only)
        load_env
        test_deployment
        ;;
    --status)
        load_env
        show_status
        ;;
    *)
        check_root
        main
        ;;
esac 