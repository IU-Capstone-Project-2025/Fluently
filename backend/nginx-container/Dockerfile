FROM nginx:alpine
COPY backend/nginx-container/nginx.conf /etc/nginx/conf.d/default.conf
COPY backend/docs /usr/share/nginx/html/nginx-container
COPY frontend-website /usr/share/nginx/html/nginx-container/index.html
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]