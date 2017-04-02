#!/bin/bash

# $1: Droplet id's (CSV)

curl -sS -X POST "https://api.digitalocean.com/v2/load_balancers/dfc8c5be-a4ae-4fdb-9e70-6c4a11999e59/droplets" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $DOTOKEN" \
    -d '{"droplet_ids": ['$1']}' 
 