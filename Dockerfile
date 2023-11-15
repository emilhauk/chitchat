FROM golang:1.21-alpine

LABEL authors="Emil Haukeland <emil.haukeland@gmail.com>"

ADD build/* .
ADD templates .