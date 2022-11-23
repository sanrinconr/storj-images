#!/bin/bash

# If docker compose is not gone down, when a local environment is running a collision of ports is generated.
trap graceful_shutdown SIGINT SIGQUIT SIGTERM

graceful_shutdown()
{
  echo -e "\nDown container..."
  docker-compose down
  echo "\n\n Finished \n\n"
  exit
}
docker-compose down
docker-compose build
docker-compose up -d
docker logs storj-images-backend-container -f 
