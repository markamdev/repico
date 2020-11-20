#!/bin/bash


API_PATH="/v1/gpio"

if [ -z "$API_ENDPOINT" ];
then
	export API_ENDPOINT="localhost:9000"
fi

# pre-defined list of GPIO pins used for testing
PINS=(2 3 4)
# number of trials
TRIALS=(1 2 3 4)

echo "Testing output GPIOs ${PINS[@]} at $API_ENDPOINT"

for t in ${TRIALS[@]}; do
    echo "** Trial number $t"
    for p in ${PINS[@]}; do
        echo "*** Pin number"
        curl -i -X PUT -H "ContentType:application/json" -d '{"state" :1}' $API_ENDPOINT/v1/gpio/$p
        sleep 1
        curl -i -X PUT -H "ContentType:application/json" -d '{"state" :0}' $API_ENDPOINT/v1/gpio/$p
        sleep 0
    done
done
