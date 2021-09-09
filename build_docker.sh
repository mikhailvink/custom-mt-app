#!/bin/bash

set -e
IMAGE_TAG=$1
IMAGE=registry.jetbrains.team/mau/crowdin-grazie:${IMAGE_TAG:-latest}
docker build . --tag $IMAGE
docker push $IMAGE
