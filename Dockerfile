FROM golang:1.24.6-alpine AS builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o /build/metar-provider .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /metar-provider

COPY --from=builder /build/metar-provider .

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["./metar-provider"]