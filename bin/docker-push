#! /bin/bash

export REPO=$1
docker login -u $DOCKER_USER -p $DOCKER_PASS \
  && docker push $REPO \
  && echo "docker-push success"
