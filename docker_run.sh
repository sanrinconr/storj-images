docker stop storj-images
docker rm storj-images
docker build -t storj-image .
docker run -e TOKEN_STORJ=$TOKEN_STORJ -e ENV=$ENV --name storj-images -t storj-image