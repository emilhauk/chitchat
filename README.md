# chitchat

[![Go](https://github.com/emilhauk/chitchat/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/emilhauk/chitchat/actions/workflows/go.yml)

Aims to be easy to setup for anyone interested in a privately hosted messenger-like service.

Demo user available: demo.user@example.com / test

This service is by no means complete. Follow [issues](https://github.com/emilhauk/chitchat/issues) to keep updated.

## Prerequisites
1. Docker 20.10 or greater (https://docs.docker.com/get-docker/)

## Kicking the tires
1. `docker compose up`

## Start developing
1. Golang 1.21 or greater (https://go.dev/doc/install)
2. `docker compose up -d --scale chitchat=0` to start database only
3. `go mod tidy` to install dependencies
4. `go run main.go` to start

## Troubleshooting
### Not seeing any changes after pulling repo
- If you're running everything in docker compose, you may need to trigger a manual rebuild
  1. `docker compose build`
  2. `docker compose up`