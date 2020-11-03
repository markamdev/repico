#!/bin/bash

API_PATH="/v1/gpio"

echo "GET on $API_PATH"
curl -i -X GET localhost:9000$API_PATH
sleep 3

echo "GET on $API_PATH/10"
curl -i -X GET localhost:9000$API_PATH/10
sleep 3

echo "PUT on $API_PATH/10"
curl -i -X PUT -H "ContentType:application/json" -d '{"status" :"on"}' localhost:9000/v1/gpio/10
sleep 3

echo "POST on $API_PATH"
curl -i -X POST localhost:9000$API_PATH
sleep 3

echo "POST on $API_PATH/10"
curl -i -X POST localhost:9000$API_PATH/10
sleep 3

