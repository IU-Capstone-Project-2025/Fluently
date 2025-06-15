FROM nginx:alpine
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY ./index.html /usr/share/nginx/html/index.html
EXPOSE 80
EXPOSE 443
CMD ["nginx", "-g", "daemon off;"]