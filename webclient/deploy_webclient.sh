#!/bin/bash

# Build the docker container
docker build -t austinsmart/info344webclient . --no-cache

# Push the docker container
docker push austinsmart/info344webclient:latest

# SSH
ssh austin@138.68.41.64 ' 
sudo docker stop webclient
sudo docker rm webclient
sudo docker pull austinsmart/info344webclient:latest
sudo docker run --restart unless-stopped --name webclient -p 443:443 -p 80:80 -v /etc/info344.austinsmart.com.pem:/etc/info344.austinsmart.com.pem:ro -v /etc/info344.austinsmart.com.key:/etc/info344.austinsmart.com.key:ro -d austinsmart/info344webclient
sudo docker ps -a
'
