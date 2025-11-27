FROM golang:1.24.6-alpine AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -tags "http grpc" -o /build/metar-service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /metar-service

COPY --from=builder /build/metar-service .

ENTRYPOINT ["./metar-service"]