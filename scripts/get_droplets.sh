#!/bin/bash

curl -X -sS GET "https://api.digitalocean.com/v2/droplets?page=1&per_page=1" \
    -H "Authorization: Bearer $DOTOKEN"  \
    -H "Content-Type: application/json"
    

