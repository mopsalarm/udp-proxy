#!/bin/sh
set -e

glide install

CGO_ENABLED=0 go build -a

docker build -t mopsalarm/udp-proxy .
docker push mopsalarm/udp-proxy
