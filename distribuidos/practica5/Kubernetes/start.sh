rm Dockerfiles/servidor/servidor
rm Dockerfiles/cliente/cliente

cd raft/cmd/servidor
CGO_ENABLED=0 go build -o ../../../Dockerfiles/servidor/servidor main.go
cd ../cliente
CGO_ENABLED=0 go build -o ../../../Dockerfiles/cliente/cliente main.go


cd ../../../Dockerfiles/servidor
docker build . -t 127.0.0.1:5000/servidor:latest
docker push 127.0.0.1:5000/servidor:latest

cd ../cliente
docker build . -t 127.0.0.1:5000/cliente:latest
docker push 127.0.0.1:5000/cliente:latest


cd ../..
./go_statefulset.sh