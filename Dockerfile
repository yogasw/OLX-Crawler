FROM golang:alpine as servicewa
WORKDIR /usr/local
RUN apk add --no-cache git
ENV GO_VERSION=1.14.3
ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH /usr/local/go/bin:/go/bin:$PATH
WORKDIR /code/serviceSendTextMessage
COPY serviceSendTextMessage .
RUN ls
RUN go get .
RUN go build .

FROM php:alpine
RUN apk add --update bash
WORKDIR /code/crawl-olx
COPY crawl-olx .
RUN docker-php-ext-install sockets
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/bin --filename=composer
RUN composer install
#RUN php crawl-data-motor.php
COPY --from=servicewa /code/serviceSendTextMessage /code/serviceSendTextMessage
WORKDIR /code/
RUN mkdir /tmp/WhatsAppSession/
RUN ls /tmp/

# Run the cron
RUN echo '*/30 * * * * /usr/local/bin/php /code/crawl-olx/crawl-data-motor.php' > /etc/crontabs/root
CMD crond -l 2 -f

#RUN /code/serviceSendTextMessage/serviceSendTextMessage
CMD ["/bin/sh"]
