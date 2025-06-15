FROM nginx:alpine
COPY ../frontend-website/nginx.conf /etc/nginx/conf.d/default.conf
COPY ../frontend-website/index.html /usr/share/nginx/html/index.html
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]