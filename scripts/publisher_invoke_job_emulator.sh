#!/bin/bash

set -euC

go build -o bin/pubsub_publisher -ldflags "-w -s" ./src/pubsub_publisher
while true
do
  bin/pubsub_publisher
  sleep 60
done