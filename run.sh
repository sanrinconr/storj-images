#!/bin/bash
docker-compose down
docker-compose build
docker-compose up -d
docker logs storj-images-backend-container -f