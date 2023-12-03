docker stop kind-registry
docker rm kind-registry

docker rmi $(docker images -q) -f 

docker image prune -a 