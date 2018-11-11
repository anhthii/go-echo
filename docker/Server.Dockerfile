FROM golang:alpine as builder

WORKDIR /go/src/github.com/anhthii/go-echo

RUN apk update && \
    apk upgrade --update-cache --available

RUN apk add --no-cache git

RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/anhthii/go-echo/

# Install library dependencies
RUN dep ensure -vendor-only

COPY . /go/src/github.com/anhthii/go-echo

RUN go build -o webserver main.go


FROM alpine:latest

ADD ca-certificates.crt /etc/ssl/certs/

WORKDIR /usr/app/go-echo

COPY --from=builder /go/src/github.com/anhthii/go-echo/webserver .

EXPOSE 3000

CMD ["./webserver"]
