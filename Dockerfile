FROM golang:latest as builder

WORKDIR /build
ADD . /
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o OnlyBot ./cmd

FROM alpine
COPY --from=builder /build/OnlyBot /bin
ENTRYPOINT ./bin/OnlyBot
