#!/bin/bash

# get primary zone
zone=$(curl -sS -X GET "https://api.cloudflare.com/client/v4/zones?name=austinsmart.com&status=active&page=1&per_page=20&order=status&direction=desc&match=all" \
     -H "X-Auth-Email: $CFEMAIL" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result[0].id)

# Verify SSL for zone
result=$(curl -sS -X GET "https://api.cloudflare.com/client/v4/zones/$zone/ssl/verification" \
     -H "X-Auth-Email: $CFEMAIL" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result)

# Print
for i in "${result[@]}"
do
   echo "$i"
done
