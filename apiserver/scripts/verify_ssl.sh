#!/bin/bash

zone=$(curl -X GET "https://api.cloudflare.com/client/v4/zones?name=austinsmart.com&status=active&page=1&per_page=20&order=status&direction=desc&match=all" \
     -H "X-Auth-Email: 4.smart.austin@gmail.com" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result[0].id)

result=$(curl -X GET "https://api.cloudflare.com/client/v4/zones/$zone/ssl/verification" \
     -H "X-Auth-Email: 4.smart.austin@gmail.com" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result)

for i in "${result[@]}"
do
   echo "$i"
done
