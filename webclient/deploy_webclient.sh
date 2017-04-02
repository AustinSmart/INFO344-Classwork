#!/bin/bash

# Build the docker container
docker build -t austinsmart/info344webclient . --no-cache

# Push the docker container
docker push austinsmart/info344webclient:latest

# SSH. root = bad
ssh root@138.68.41.64 ' 
docker stop webclient
docker rm webclient
docker pull austinsmart/info344webclient:latest
docker run --name webclient -p 443:443 -p 80:80 -v /etc/info344.austinsmart.com.pem:/etc/info344.austinsmart.com.pem -v /etc/info344.austinsmart.com.key:/etc/info344.austinsmart.com.key -dt austinsmart/info344webclient
docker ps -a
'
