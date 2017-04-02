#!/bin/bash

# Get all droplets in load balancer
droplet_ids=$(sh ./get_load_balancer.sh | jq -r '.load_balancers[0].droplet_ids')
echo Droplets: ${droplet_ids}

echo Which would you like to remove?
read droplet

echo "Do you want to destroy it as well? (y/n)"
read destroy

# Removed selected droplet
if([ "$droplet" != "" ]) then
    curl -X DELETE "https://api.digitalocean.com/v2/load_balancers/dfc8c5be-a4ae-4fdb-9e70-6c4a11999e59/droplets" \
        -H "Authorization: Bearer $DOTOKEN" \
        -H "Content-Type: application/json" \
        -d '{"droplet_ids": ['$droplet']}' 
else
    echo Please enter a droplet id
fi    

# Destroy droplet if 'y'
if([ "$destroy" == "y" ]  || [ "$destroy" == "yes" ]) then
    curl -X DELETE "https://api.digitalocean.com/v2/droplets/$droplet" \
        -H "Authorization: Bearer $DOTOKEN" \
        -H "Content-Type: application/json" 
fi
