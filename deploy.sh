#!/bin/bash

set -e

# Config
IMAGE_NAME="assignment2-service"
CONTAINER_NAME="assignment2-container"
FIREBASE_KEY_PATH="$(pwd)/credentials/firebase-key.json"
CONTAINER_KEY_PATH="/credentials/firebase-key.json"
PORT=8080

echo "Building Docker image..."
docker build -t $IMAGE_NAME .

echo "Stopping and removing old container (if any)..."
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

echo "Running new container..."
docker run -d \
  --name $CONTAINER_NAME \
  -p $PORT:$PORT \
  -v "$FIREBASE_KEY_PATH":"$CONTAINER_KEY_PATH" \
  $IMAGE_NAME

echo "Deployment complete"
echo "Container logs: docker logs -f $CONTAINER_NAME"
echo "Visit: http://<your-vm-ip>:$PORT/dashboard/v1/registrations"
