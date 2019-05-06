# Builder
FROM golang:1.10-alpine as builder

LABEL maintainer "klautcomputing"
WORKDIR /go/src/github.com/Shopify/voucher
COPY . .
RUN apk --no-cache add \
    git \
    make && \
    make all

# Final build
FROM alpine:3.8

COPY --from=builder /go/src/github.com/Shopify/voucher/build/voucher_server /usr/local/bin/voucher_server
COPY --from=builder /go/src/github.com/Shopify/voucher/entrypoint.sh /usr/local/entrypoint.sh
COPY --from=builder /go/src/github.com/Shopify/voucher/config/config.toml /etc/voucher/config.toml
COPY config/secrets.production.ejson /etc/voucher/secrets.production.ejson

RUN apk add --no-cache \
    ca-certificates && \
    addgroup -S -g 10000 voucher && \
    adduser -S -u 10000 -G voucher voucher

USER 10000:10000

ENTRYPOINT ["/usr/local/entrypoint.sh"]
