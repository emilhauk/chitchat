# chitchat

Aims to be easy to setup for anyone interested in a privately hosted messenger-like service.
It's just about unusable in its current state, but feel free to contribute.

## Prerequisites
1. Docker 20.10 or greater (https://docs.docker.com/get-docker/)

## Kicking the tires
1. `docker compose up`

## Start developing
1. Golang 1.21 or greater (https://go.dev/doc/install)
2. `docker compose up -d --scale chitchat=0` to start database only
3. `go mod tidy` to install dependencies
4. `go run main.go` to start

Demo user available: demo.user@example.com / test