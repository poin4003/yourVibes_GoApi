FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

ARG BUILD_ENVIRONMENT=prod
ENV YOURVIBES_SERVER_CONFIG_FILE=$BUILD_ENVIRONMENT

COPY . .

RUN go mod download
RUN go build -o yourvibes_api_server ./cmd/server

FROM alpine:latest

RUN apk update && apk add --no-cache \
    ffmpeg

COPY ./config /config
COPY ./templates /templates

COPY --from=builder /build/yourvibes_api_server /

ENTRYPOINT [ "/yourvibes_api_server", "config/prod.yaml" ]
