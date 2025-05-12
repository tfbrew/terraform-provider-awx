#! /bin/bash

# echo "AAP endpoint"
# curl --request GET --url https://aap-provider-test.morford.dev/api/ | jq


echo "AAP ping"
curl --request GET --url https://aap-provider-test.morford.dev/api/gateway/v1/ping/ | jq


# echo "-----------------"

# echo "AWX endpoint"
# curl --request GET --url https://awx-provider-test.morford.dev/api/ | jq

# echo ""

# curl --request GET --url https://aap-provider-test.morford.dev/api/controller/v2/ping/ --header 'Content-Type: application/json'


