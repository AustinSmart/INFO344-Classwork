#!/bin/bash

# Build the Go binary
GOOS=linux go build

# Build the docker container
docker build -t austinsmart/info344apiserver . --no-cache

# Push the docker container
docker push austinsmart/info344apiserver:latest

# SSH. root = bad
ssh root@138.197.219.125 ' 
docker stop apiserver
docker rm apiserver
docker pull austinsmart/info344apiserver:latest
docker run --name apiserver -p 443:443 -v /etc/info344api.austinsmart.com.pem:/etc/info344api.austinsmart.com.pem -v /etc/info344api.austinsmart.com.key:/etc/info344api.austinsmart.com.key -e TLSCERT=/etc/info344api.austinsmart.com.pem -e TLSKEY=/etc/info344api.austinsmart.com.key -d austinsmart/info344apiserver
docker ps
'
