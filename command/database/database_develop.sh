docker build -t boilerplate-db ./
docker run -d --name boilerplate-db -p 5436:5432 boilerplate-db