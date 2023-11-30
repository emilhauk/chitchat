# chitchat

[![Go](https://github.com/emilhauk/chitchat/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/emilhauk/chitchat/actions/workflows/go.yml)

Aims to be easy to set up for anyone interested in a privately hosted messenger-like service.

Demo user available: demo.user@example.com / test

This service is by no means complete. Follow [issues](https://github.com/emilhauk/chitchat/issues) to keep updated.

## Prerequisites
1. Docker 20.10 or greater (https://docs.docker.com/get-docker/).

## Kicking the tires
1. `docker compose up`

## Start developing
1. Golang 1.21 or greater (https://go.dev/doc/install).
2. `docker compose up -d --scale chitchat=0` to start database only.
3. `go mod tidy` to install dependencies.
4. `go run main.go` to start.

## Running migrations
We only have one migration yet. The initial one, but when we reach V1 there will need to be separate ones. The current migration strategy is:
1. `docker compose exec db mysql -uchitchat -ppassword chitchat` (please use something difrferent in production).
   - PS: The db container need to be fully running before attempting to connect. Ut has a health check that should read "healtly" in `docker compose ps`
2. Paste and run new migrations manually from the schema folder.

## Troubleshooting
### Not seeing any changes after pulling repo
- If you're running everything in docker compose, you may need to trigger a manual rebuild:
  1. `docker compose build`
  2. `docker compose up`
### If chitchat fails to start check logs
- It logs complains about prepared statements you'll need to either manually run migrations, or just nuke the database entierly. The latter will of course be unacceptable in production :)
  1. `docker compose down` (nukes database)
  2. `docker compose up`

## Instabilities be damned, lets run it now anyway!
Very well then! 
You will need a hosting service, or host it at home. There are plenty of guides on now to host web services, so I'll leave it up to you. But here's the short list of requirements:
 - MySQL / MariaDB
 - Able to run containers
   - Although you could build it manually and run it directly if you'd like
 - A domain (preferrably) and SSL cert to get https.
   - At this time you _may_ run on http, but I encourage you to read up on why that's perhaps not the greatest idea.
   - In addition, you won't be able to make use of future services like; web push, passkey login etc.

I encourage you to have a look at the included docker-compose.yml file to check out available environment variables to set.