#!/bin/bash

#curl -v -XPOST http://127.0.0.1:4082/v0/tools/vegeta -d '{}'
curl -v -XPOST http://127.0.0.1:4082/v0/tools/vegeta -d '{
   "target":    "GET http://127.0.0.1:8080",
   "duration":  5, 
   "timeout":   2,
   "rate":      1000, 
   "workers":    4,
   "connections": 100,
   "statsd": {
       "enable": true,
       "host":  "127.0.0.1",
       "port":  8125,
       "prefix": "vegeta"
   }
}'
