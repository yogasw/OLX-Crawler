FROM golang:1.14.3-alpine3.11

# Set the working directory to /app
WORKDIR /app

# COPY requirements to /app dir
COPY . .
RUN ls
RUN go
RUN echo "hello world"

