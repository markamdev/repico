#!/bin/bash

API_PATH="/v1/gpio/number"

if [ -z "$API_ENDPOINT" ];
then
	export API_ENDPOINT="localhost:9000"
fi

if [ -z "$PIN_NUMBER" ];
then
	export PIN_NUMBER="4"
fi

echo "GET on $API_PATH"
curl -i -X GET $API_ENDPOINT$API_PATH
sleep 3

echo "GET on $API_PATH/$PIN_NUMBER"
curl -i -X GET $API_ENDPOINT$API_PATH/$PIN_NUMBER
sleep 3

echo "PUT '1' on $API_PATH/$PIN_NUMBER"
curl -i -X PUT -H "ContentType:application/json" -d '{"state" :1}' $API_ENDPOINT$API_PATH/$PIN_NUMBER
sleep 3

echo "PUT '0' on $API_PATH/$PIN_NUMBER"
curl -i -X PUT -H "ContentType:application/json" -d '{"state" :0}' $API_ENDPOINT$API_PATH/$PIN_NUMBER
sleep 3

echo "POST on $API_PATH"
curl -i -X POST $API_ENDPOINT$API_PATH
sleep 3

echo "POST on $API_PATH/$PIN_NUMBER"
curl -i -X POST $API_ENDPOINT$API_PATH/$PIN_NUMBER
sleep 3
