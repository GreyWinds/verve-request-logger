
# Steps to Run:
docker network create verve-net

docker build -t verve .

docker run -p 6379:6379 --name redis --network verve-net redis

docker run -p 8080:8080 --network verve-net verve


# To view the log file in a stream:
docker exec -it *<container-**name**-created>* /bin/sh

tail -f unique_requests.log
