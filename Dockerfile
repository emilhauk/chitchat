FROM golang:1.21-alpine AS build

# git is required for fetching dependencies
RUN apk update && apk add git make && mkdir -p /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    APP_VERSION=${APP_VERSION}

WORKDIR /app

# download modules first to be able to cache the resulting layer
ADD go.mod go.sum ./
RUN go mod download

ADD . .

RUN make release

FROM alpine:latest AS release

LABEL authors="Emil Haukeland <emil.haukeland@gmail.com>"

ENV LC_ALL=en_US.UTF-8
ENV LC_LANG=en_US.UTF-8
ENV LC_LANGUAGE=en_US.UTF-8

RUN mkdir -p /app/{templates,static} && \
    apk update && \
    apk add tzdata ca-certificates curl && \
    update-ca-certificates 2>/dev/null || true && \
    mkdir -p /app/bin

ENV TZ Europe/Oslo

COPY --from=build /app/build/chitchat /app/templates /app/
COPY --from=build /app/templates /app/templates/
COPY --from=build /app/static /app/static/
RUN chmod +x /app/chitchat

WORKDIR /app

RUN addgroup -S chitchat && adduser -S chitchat -G chitchat
USER chitchat

ENTRYPOINT ["/app/chitchat"]

EXPOSE ${PORT:-3333}
