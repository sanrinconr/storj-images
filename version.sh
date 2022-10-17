# Generate a new image version and push to the registry

#docker-compose  build
imageName=storj-images
actualImage=$imageName-backend
TAG=${REGISTRY}/$imageName:$1

# Docker compose dont allow custom tag names
docker-compose build

## Rename to prepare to upload to the registry
docker tag $imageName-backend $TAG
docker push $TAG

# Remove bad named imaged
docker rmi $actualImage
