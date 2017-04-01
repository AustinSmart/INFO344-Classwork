#!/bin/bash

# $1: Load balancer id

curl -X GET "https://api.digitalocean.com/v2/load_balancers/$1" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" 
