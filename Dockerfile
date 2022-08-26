FROM golang:latest

WORKDIR /bot
COPY go.mod ./
COPY go.sum ./

RUN apt-get update

COPY . .
RUN go mod download

RUN go build -o OnlyBot ./cmd
ENTRYPOINT ./OnlyBot