FROM golang:1.17 AS builder

COPY . /build
WORKDIR /build

RUN GOOS=linux go build -v -o app


FROM debian:buster-slim
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
RUN mkdir /app && adduser -D -H -u 1000 -s /bin/false app
COPY docker-entrypoint.sh /docker-entrypoint.sh
COPY --from=builder --chown=1000:1000 /build/app /app/priceupdater
RUN chown -R 1000:1000 /app

USER 1000

ENTRYPOINT ["/docker-entrypoint.sh", "/app/priceupdater"]
