# Generate a new image version and push to the registry

#docker-compose  build
imageName=${REGISTRY}/storj-images-backend
TAG=$imageName:$1

# Docker compose dont allow custom tag names
docker-compose build

## Rename to prepare to upload semantic version to the registry
docker tag $imageName $TAG
docker push $TAG

git tag $1
# Remove bad named imaged
docker rmi $imageName

