#!/bin/sh

echo "Building image ${IMAGE_PATH}"

# build docker image
docker build -t ${IMAGE_PATH} .
gcloud docker -- push ${IMAGE_PATH}