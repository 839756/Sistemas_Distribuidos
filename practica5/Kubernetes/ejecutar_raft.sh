
rm Dockerfiles/servidor/servidor
rm Dockerfiles/cliente/cliente

cd raft/cmd/servidor
CGO_ENABLE=0 go build -o ../../../Dockerfiles/servidor/servidor main.go

cd ../cliente
CGO_ENABLE=0 go build -o ../../../Dockerfiles/cliente/cliente main.go

cd ../../../Dockerfiles/cliente
sudo docker build . -t localhost:5001/cliente:latest
sudo docker push localhost:5001/cliente:latest

cd ../servidor
sudo docker build . -t localhost:5001/servidor:latest
sudo docker push localhost:5001/servidor:latest

cd ../..
sudo ./go_statefulset.sh
