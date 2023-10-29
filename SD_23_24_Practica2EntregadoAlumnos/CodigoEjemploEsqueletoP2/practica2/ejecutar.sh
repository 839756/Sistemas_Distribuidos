#!/bin/bash

# Definir usuario
USER="a839756"

# Iniciar el programa barrera.go en la máquina 192.168.3.9:29121
ssh -n $USER@192.168.3.9 "cd ~/practica2/; go run barrera.go" &

# Darle tiempo al programa barrera.go para que se inicie correctamente
sleep 2

# Array de IPS
IPS=("192.168.3.9" "192.168.3.10" "192.168.3.11" "192.168.3.12")

# Iniciar el programa lector en 192.168.3.9:29120
ssh -n "$USER@${IPS[0]}" "cd ~/practica2/cmd/lector; go run main.go 1" &

# Iniciar el programa lector en 192.168.3.10:29120
ssh -n "$USER@${IPS[1]}" "cd ~/practica2/cmd/lector; go run main.go 2" &

# Iniciar el programa escritor en 192.168.3.11:29120
ssh -n "$USER@${IPS[2]}" "cd ~/practica2/cmd/escritor; go run main.go 3" &

# Iniciar el programa escritor en 192.168.3.12:29120
ssh -n "$USER@${IPS[3]}" "cd ~/practica2/cmd/escritor; go run main.go 4" &