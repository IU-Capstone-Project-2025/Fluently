FROM nginx:alpine
COPY backend/nginx-container/nginx.conf /etc/nginx/conf.d/default.conf
COPY backend/nginx-container/nginx-headers.conf /etc/nginx/nginx-headers.conf
COPY backend/docs /usr/share/nginx/html/nginx-container
COPY frontend-website /usr/share/nginx/html/main-site/
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]