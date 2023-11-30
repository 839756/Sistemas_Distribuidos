#!/bin/bash

kind delete cluster
./kind-with-registry.sh

rm Dockerfiles/servidor/servidor
rm Dockerfiles/cliente/cliente

cd raft/cmd/servidor
CGO_ENABLE=0 go build -o ../../../Dockerfiles/servidor/servidor main.go

cd ../cliente
CGO_ENABLE=0 go build -o ../../../Dockerfiles/cliente/cliente main.go

cd ../../../Dockerfiles/cliente
# Build Docker image for servidor and push it
docker build -t localhost:5000/servidor:latest .
docker push localhost:5000/servidor:latest

cd ../cliente
# Build Docker image for cliente and push it
docker build -t localhost:5000/cliente:latest .
docker push localhost:5000/cliente:latest

cd ../..
kubectl create -f statefulset_go.yaml