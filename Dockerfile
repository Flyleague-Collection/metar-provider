FROM golang:1.25.5-alpine AS builder

WORKDIR /build

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -tags "http telemetry" -o /build/metar-service .

FROM alpine:latest AS runtime

RUN apk update
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl

WORKDIR /service

COPY --from=builder /build/metar-service .

ENTRYPOINT ["./metar-service"]