#!/bin/bash

# $1: VM name
# $2: docker container name - optional
# $3: docker image name - optional
# $@: -lb will add to loadbalancer - optional

# Create a new Ubuntu 16.04 x64 VM in the SFO1 region with 512mb of RAM named $1
if ([ -z "$1" ]) then
    droplet_id=$(curl -sS  -X POST "https://api.digitalocean.com/v2/droplets" \
        -H "Authorization: Bearer $DOTOKEN" \
        -H "Content-Type: application/json" \
        -d'{"name":"'$1'","region":"sfo2","size":"512mb","image":"ubuntu-16-04-x64",

        "user_data":"
    #cloud-config

    packages:
        - docker.io
    runcmd:
        - docker run -d --name '$2' -p 80:80 austinsmart/'$3':latest
    "}' | jq -r '.droplet.id')
else # No "docker run" command
    droplet_id=$(curl -sS  -X POST "https://api.digitalocean.com/v2/droplets" \
        -H "Authorization: Bearer $DOTOKEN" \
        -H "Content-Type: application/json" \
        -d'{"name":"'$1'","region":"sfo2","size":"512mb","image":"ubuntu-16-04-x64",

        "user_data":"
    #cloud-config

    packages:
        - docker.io
    "}' | jq -r '.droplet.id')
fi

echo Droplet created: $droplet_id

# Check if adding to load balancer
addlb=false

while [ "$1" != "" ]; do
    case $1 in
        -lb ) addlb=true
    esac
    shift
done

if ([ $addlb == true ]) then

    # Wait until the droplet status is "active" and then add it to the load balancer
    status=""
    while [ "$status" != "active" ]
    do
        status=$(curl -sS  -X GET  "https://api.digitalocean.com/v2/droplets/$droplet_id" \
            -H "Authorization: Bearer $DOTOKEN" \
            -H "Content-Type: application/json" | jq -r '.droplet.status')

        echo "Waiting for Droplet to come online"
        
        sleep 10
    done

    sh ./add_droplet_to_load_balancer.sh $droplet_id
fi




