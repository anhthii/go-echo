FROM node:alpine as builder

RUN apk update && \
    apk upgrade --update-cache --available

RUN apk add git

RUN git clone https://github.com/anhthii/Echo /usr/app/Echo

WORKDIR /usr/app/Echo

RUN npm install

RUN cd app/constant && \
    sed -i 's/localhost/server/g' endpoint_constant.js

RUN npm run webpack:prod

FROM alpine:latest

RUN apk update && \
    apk upgrade --update-cache --available

RUN apk add nginx

RUN mkdir -p /run/nginx && \
    mkdir /usr/share/nginx && \
    mkdir /usr/share/nginx/html

COPY --from=builder /usr/app/Echo/public /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/

EXPOSE 8080

ENTRYPOINT ["nginx", "-g", "daemon off;"]



