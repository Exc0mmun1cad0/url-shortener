# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23.4 AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /url-shortener cmd/url-shortener/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /migrator cmd/migrator/main.go

# TODO: add running tests hete

# Deploy the application in distroless image
FROM alpine:latest AS release

ARG CONFIG_PATH
ENV CONFIG_PATH=${CONFIG_PATH}

ARG MIGRATIONS_PATH
ENV MIGRATIONS_PATH=${MIGRATIONS_PATH}

WORKDIR /

COPY --from=build /url-shortener /url-shortener
COPY --from=build /migrator /migrator
COPY .env /
COPY ${CONFIG_PATH} /config.yaml
COPY ${MIGRATIONS_PATH} /


