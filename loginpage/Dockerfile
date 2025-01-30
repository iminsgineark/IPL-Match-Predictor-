FROM golang:1.23-alpine AS builder

LABEL maintainer="Utkrist Ark"
Label description="Dockerfile for go login app"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]

