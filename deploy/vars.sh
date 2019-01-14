#!/bin/sh
PROJECT_ID=epigos-ai
BRANCH_NAME=$(git branch | sed -n '/\* /s///p')
ENV=$1
PROD="prod"
DEV="dev"

VERSION=$BRANCH_NAME

if [[ $ENV == $PROD ]]; then
    source ./deploy/tag.sh
    TAG=$(git describe --tags --match=${VERSION}* --abbrev=0)
    VERSION=$TAG
else
    PROJECT_ID=${PROJECT_ID}-dev
fi

IMAGE_PATH=eu.gcr.io/epigos-ai-${ENV}/newsbot:${VERSION}

echo "Setting env vars"

export IMAGE_PATH=$IMAGE_PATH
export VERSION=$VERSION
export PROJECT_ID=$PROJECT_ID
export ENV=$ENV

# configure gcloud
gcloud config set account ${GCLOUD_ACCOUNT}
gcloud config set project ${PROJECT_ID}