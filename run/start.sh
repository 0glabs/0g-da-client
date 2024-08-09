#!/bin/bash
set -e

IMAGE_NAME="0gclient:latest"
CONTAINER_NAME="daclient"

HOST_PORT=51001
CONTAINER_PORT=51001

docker run -d -v .:/runtime --name=$CONTAINER_NAME -p $HOST_PORT:$CONTAINER_PORT --restart=always $IMAGE_NAME

echo "Checking if container is running..."
docker ps --filter "name=$CONTAINER_NAME"
