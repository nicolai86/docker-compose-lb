#!/usr/bin/env bash

set -e
set -u

[[ -e reverse-proxy ]] && rm reverse-proxy
GOOS=linux GOARCH=amd64 go build -o reverse-proxy main.go
chmod +x reverse-proxy

docker build -t nicolai86/docker-compose-reverse-proxy .
