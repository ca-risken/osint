#!/bin/bash -e

cd "$(dirname "$0")"

# load env
. ../env.sh

# setting remote repository
TAG="local-test-$(date '+%Y%m%d')"
IMAGE_OSINT="osint/osint"
IMAGE_SUBDOMAIN="osint/subdomain"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
REGISTORY="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

# build & push
aws ecr get-login-password --region ${AWS_REGION} \
  | docker login \
    --username AWS \
    --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com

docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_OSINT}:${TAG} ../src/osint/
docker build --build-arg GITHUB_USER=${GITHUB_USER} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} -t ${IMAGE_SUBDOMAIN}:${TAG} ../src/subdomain/

docker tag ${IMAGE_OSINT}:${TAG}     ${REGISTORY}/${IMAGE_OSINT}:${TAG}
docker tag ${IMAGE_SUBDOMAIN}:${TAG} ${REGISTORY}/${IMAGE_SUBDOMAIN}:${TAG}

docker push ${REGISTORY}/${IMAGE_OSINT}:${TAG}
docker push ${REGISTORY}/${IMAGE_SUBDOMAIN}:${TAG}
