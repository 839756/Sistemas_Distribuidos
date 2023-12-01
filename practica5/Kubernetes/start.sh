#!/bin/bash

kind delete cluster
docker stop kind-registry
docker rm kind-registry
./kind-with-registry.sh

rm Dockerfiles/servidor/servidor
rm Dockerfiles/cliente/cliente

cd raft/cmd/servidor
CGO_ENABLE=0 go build -o ../../../Dockerfiles/servidor/servidor main.go

cd ../cliente
CGO_ENABLE=0 go build -o ../../../Dockerfiles/cliente/cliente main.go

cd ../../../Dockerfiles/servidor
# Build Docker image for servidor and push it
docker build -t 127.0.0.1:5000/servidor:latest .
docker push 127.0.0.1:5000/servidor:latest

# docker tag localhost:5001/servidor:latest 839756/sistemas-sistribuidos:servidor
# docker push 839756/sistemas-sistribuidos:servidor

cd ../cliente
# Build Docker image for cliente and push it
docker build -t 127.0.0.1:5000/cliente:latest .
docker push 127.0.0.1:5000/cliente:latest

# docker tag localhost:5001/cliente:latest 839756/sistemas-sistribuidos:cliente
# docker push 839756/sistemas-sistribuidos:cliente

cd ../..
kubectl create -f statefulset_go.yaml