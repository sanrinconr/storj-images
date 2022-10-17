#!/bin/bash
# Run a version already pushed into the registry, you can simulate the prod environment in local :)
lastStable=1.0.0
imageName=storj-images
container=$imageName-container

version=$1
if [ -z "$1" ]
  then
    echo "No version passed, using the last stable: $lastStable"
    version=$lastStable
fi
tag=${REGISTRY}/$imageName:$version

ENV="prod"
PORT=6002

docker stop $container
docker rm $container
docker pull $tag
sleep 2

docker run -p=$PORT:$PORT -e ENV=$ENV -e PORT=$PORT --name $container -t $tag
