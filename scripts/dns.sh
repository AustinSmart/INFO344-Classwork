#!/bin/bash

# $1: Record type

# View DNS Records for austinsmart.com, returns all records if no type is supplied

zone=$(curl -sS -X GET "https://api.cloudflare.com/client/v4/zones?name=austinsmart.com&status=active&page=1&per_page=20&order=status&direction=desc&match=all" \
     -H "X-Auth-Email: $CFEMAIL" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result[0].id)

# Get all
if ([ -z "$1" ]) then
dns=$(curl -sS -X GET "https://api.cloudflare.com/client/v4/zones/$zone/dns_records?page=1&per_page=20&order=type&direction=desc&match=all" \
     -H "X-Auth-Email: $CFEMAIL" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result)
else
# Get type
dns=$(curl -sS -X GET "https://api.cloudflare.com/client/v4/zones/$zone/dns_records?type=$1&page=1&per_page=20&order=type&direction=desc&match=all" \
     -H "X-Auth-Email: $CFEMAIL" \
     -H "X-Auth-Key: $CFTOKEN" \
     -H "Content-Type: application/json" | jq -r .result)
fi

# Print
for i in "${dns[@]}"
do
   echo "$i"
done