#!/bin/bash

curl -X GET "https://api.digitalocean.com/v2/droplets?page=1&per_page=1" \
    -H "Authorization: Bearer $TOKEN"  \
    -H "Content-Type: application/json"
    

