FROM nginx:alpine
COPY backend/docs /usr/share/nginx/html/nginx-container
COPY frontend-website /usr/share/nginx/html/main-site/
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]