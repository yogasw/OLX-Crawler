FROM golang:alpine
WORKDIR /usr/local
RUN apk update
RUN apk add git
RUN apk add --update bash
ENV GO_VERSION=1.14.3
ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH /usr/local/go/bin:/go/bin:$PATH
WORKDIR /code
COPY . .
WORKDIR serviceSendTextMessage
RUN go get .
RUN go build .



CMD ["/bin/sh"]