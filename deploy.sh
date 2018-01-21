#!/bin/bash

PROJECT_ID=$1
IMAGE=$2
VERSION=$3

echo "Deploying ${IMAGE} v${VERSION} to GCP project ${PROJECT_ID}"

docker build -t gcr.io/${PROJECT_ID}/${IMAGE}:v${VERSION} .
gcloud docker -- push gcr.io/${PROJECT_ID}/${IMAGE}:v${VERSION}
kubectl set image deployment/${IMAGE}-deploy ${IMAGE}-deploy=gcr.io/${PROJECT_ID}/${IMAGE}:v${VERSION}
