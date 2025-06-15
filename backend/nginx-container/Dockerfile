FROM nginx:alpine
COPY nginx-container/nginx.conf /etc/nginx/conf.d/default.conf
COPY docs /usr/share/nginx/html/nginx-container
COPY ../frontend-website/index.html /usr/share/nginx/html/index.html
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]