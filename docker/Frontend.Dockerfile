FROM node:alpine as builder

RUN apk update && \
    apk upgrade --update-cache --available

RUN apk add git

RUN git clone https://github.com/anhthii/Echo /usr/app/Echo

WORKDIR /usr/app/Echo

RUN npm install

RUN cd app/constant && \
    sed -i 's/localhost/server/g' endpoint_constant.js

RUN npm run build

FROM alpine:latest

RUN apk update && \
    apk upgrade --update-cache --available

RUN apk add --no-cache bash nginx 


RUN mkdir -p /run/nginx && \
    mkdir /usr/share/nginx && \
    mkdir /usr/share/nginx/html

COPY --from=builder /usr/app/Echo/public /usr/share/nginx/html

COPY ./docker/nginx.conf /etc/nginx/conf.d/

COPY ./docker/wait-for-it.sh /wait-for-it.sh

RUN chmod +x /wait-for-it.sh

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]



