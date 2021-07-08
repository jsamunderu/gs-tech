#!/bin/bash

curl -X POST http://localhost:8080/login -d '{"Username":"jamess", "Password":"testpassword"}' -v
curl -X POST http://localhost:8080/query -d '{"ProfileId":"jamess", "Token":"0f97b309-bc52-4d50-b50e-73a079b1c82b"}' -v
curl -X POST http://localhost:8080/profileDelete -d '{"ProfileId":"jamess", "Token":"0f97b309-bc52-4d50-b50e-73a079b1c82b"}' -v
curl -X POST http://localhost:8080/profileUpdate -d '{"Profile":{"ProfileId":"jamess","Name":"James","Age":33,"FavoriteColor":"Navy Blue","FavoriteOperatingSystem":"linux"},"token":"0f97b309-bc52-4d50-b50e-73a079b1c82b"}' -v
