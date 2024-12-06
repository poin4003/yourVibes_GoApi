FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o yourvibes_api_server ./cmd/server

FROM scratch

COPY ./config /config

COPY --from=builder /build/yourvibes_api_server /

ENTRYPOINT [ "/yourvibes_api_server", "config/local.yaml" ]