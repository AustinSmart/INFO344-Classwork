#!/bin/bash

# $1: Load balancer name

curl -sS -X POST  "https://api.digitalocean.com/v2/load_balancers" \
    -H "Authorization: Bearer $DOTOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name": "'$1'",
        "region": "sfo2",
        "forwarding_rules":
            [{"entry_protocol":"http",
            "entry_port":80,
            "target_protocol":"http",
            "target_port":80,
            "certificate_id":"",
            "tls_passthrough":false},

            {"entry_protocol": "https",
            "entry_port": 444,
            "target_protocol": "https",
            "target_port": 443,
            "tls_passthrough": true}
            ], 
        "health_check":
            {"protocol":"http",
            "port":80,
            "path":"/",
            "check_interval_seconds":10,
            "response_timeout_seconds":5,
            "healthy_threshold":5,
            "unhealthy_threshold":3},
            "sticky_sessions":{"type":"none"}
            }' 