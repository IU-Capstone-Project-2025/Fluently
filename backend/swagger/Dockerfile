FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY docs /usr/share/nginx/html/swagger
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]