#!/bin/bash

hashed_password=$(go run ./cmd/pwd_tool/... testpassword)
echo "SET jamess ${hashed_password}" | redis-cli -n 3

echo "GET jamess" | redis-cli -n 3

echo "SET jamess '{\"ProfileId\":\"jamess\",\"Name\":\"James\",\"Age\":33,\"FavoriteColor\":\"Navy Blue\",\"FavoriteOperatingSystem\":\"linux\"}'" | redis-cli -n 1

echo "GET jamess" | redis-cli -n 1
