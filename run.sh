#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./bin/server -port=8001 &
./bin/server -port=8002 &
./bin/server -port=8003 -api=1 &

sleep 2
echo ">>> start test"
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &
curl "http://localhost:9999/api?key=Tom" &

wait
