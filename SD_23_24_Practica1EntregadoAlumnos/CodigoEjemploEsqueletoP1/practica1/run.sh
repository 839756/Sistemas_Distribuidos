#!/bin/bash

# Ejecutar el primer script
./build_all.sh

# Ejecutar el servidor en una nueva terminal
gnome-terminal -- go run cmd/server-draft/concurrentePeticion.go

# Esperar un momento antes de ejecutar el cliente
sleep 2

# Ejecutar el cliente en otra nueva terminal
gnome-terminal -- go run cmd/client/main.go
