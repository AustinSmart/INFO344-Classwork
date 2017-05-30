#!/bin/bash

# Build the Go binary
GOOS=linux go build

# Build the docker container
docker build -t austinsmart/info344apiserver .

# Push the docker container
docker push austinsmart/info344apiserver:latest

# SSH
ssh austin@138.197.219.125 ' 
sudo docker stop apiserver
sudo docker rm apiserver
sudo docker pull austinsmart/info344apiserver:latest
sudo docker run --restart unless-stopped --name apiserver -p 443:443 \
--network api-server-net \
-v /etc/info344api.austinsmart.com.pem:/etc/info344api.austinsmart.com.pem:ro \
-v /etc/info344api.austinsmart.com.key:/etc/info344api.austinsmart.com.key:ro \
-e TLSCERT=/etc/info344api.austinsmart.com.pem \
-e TLSKEY=/etc/info344api.austinsmart.com.key \
-e DBADDR="mongo" \
-e REDISADDR="redis:6379" \
-e SESSIONKEY="khqfkhthrqbfniwhieferjduesHF" \
-e CHATBOTADDR=tautBot \
-d austinsmart/info344apiserver
sudo docker ps -a
'

