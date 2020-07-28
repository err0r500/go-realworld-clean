#!/usr/bin/env bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

BINARY=go-realworld-clean

GO111MODULE=on go build -o $(BINARY)
sudo docker-compose up --build
