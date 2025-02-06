FROM nginx:alpine

COPY cert/server.crt /etc/nginx/cert/server.crt

COPY cert/server.key /etc/nginx/cert/server.key

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80 443

