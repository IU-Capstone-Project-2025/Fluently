FROM nginx:alpine

# Install gettext for envsubst
RUN apk add --no-cache gettext

# Copy website content
COPY frontend-website /usr/share/nginx/html/main-site/

# Copy the NGINX configuration template
COPY backend/nginx-container/nginx.conf.template /etc/nginx/templates/default.conf.template

# Create a startup script that processes the template
RUN echo '#!/bin/sh' > /docker-entrypoint.d/40-process-template.sh && \
    echo 'set -e' >> /docker-entrypoint.d/40-process-template.sh && \
    echo '' >> /docker-entrypoint.d/40-process-template.sh && \
    echo '# Set default values for template variables if not provided' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'export DOMAIN=${DOMAIN:-localhost}' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'export CERT_NAME=${CERT_NAME:-localhost}' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'export webhook_secret=${WEBHOOK_SECRET:-}' >> /docker-entrypoint.d/40-process-template.sh && \
    echo '' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'echo "Processing NGINX template with variables:"' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'echo "  DOMAIN: $DOMAIN"' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'echo "  CERT_NAME: $CERT_NAME"' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'echo "  WEBHOOK_SECRET: [hidden]"' >> /docker-entrypoint.d/40-process-template.sh && \
    echo '' >> /docker-entrypoint.d/40-process-template.sh && \
    echo '# Process the template and output to the correct location' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'envsubst '\''${DOMAIN} ${CERT_NAME} ${webhook_secret}'\'' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf' >> /docker-entrypoint.d/40-process-template.sh && \
    echo '' >> /docker-entrypoint.d/40-process-template.sh && \
    echo 'echo "NGINX configuration generated successfully"' >> /docker-entrypoint.d/40-process-template.sh

# Make the script executable
RUN chmod +x /docker-entrypoint.d/40-process-template.sh

EXPOSE 80
EXPOSE 443

CMD ["nginx", "-g", "daemon off;"]