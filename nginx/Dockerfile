FROM openresty/openresty:alpine

COPY nginx.conf /usr/local/openresty/nginx/conf/nginx.conf

COPY mime.types etc/nginx/mime.types

EXPOSE 80

CMD ["openresty", "-g", "daemon off;"]

