#!/bin/bash

# Build the docker container
docker build -t austinsmart/info344chatbot .

# Push the docker container
docker push austinsmart/info344chatbot:latest

# SSH
ssh austin@138.197.219.125 ' 
sudo docker stop tautBot
sudo docker rm tautBot
sudo docker pull austinsmart/info344chatbot:latest
sudo docker run \
--name tautBot \
--network api-server-net \
-e WITAITOKEN=LL42BS3DBBDR3OKJHENQQT6CK2USTJP4 \
-d austinsmart/info344chatbot:latest 
sudo docker ps -a
'

