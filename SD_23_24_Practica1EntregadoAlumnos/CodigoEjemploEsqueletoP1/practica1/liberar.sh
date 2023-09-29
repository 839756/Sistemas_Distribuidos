#!/bin/bash

# Define las direcciones IP de las máquinas remotas
remote_ips=("192.168.3.10" "192.168.3.11" "192.168.3.12")

# Define el usuario
remote_user="$USER"

# Define el puerto que deseas verificar y matar
port_to_kill=29120

# Recorre las direcciones IP y ejecuta el comando para matar el proceso si es necesario
for ip in "${remote_ips[@]}"; do
  # Utiliza SSH para conectarte a la máquina remota y ejecutar el comando
  ssh "$remote_user@$ip" "pkill -f :$port_to_kill"
  echo "Proceso en $ip en el puerto $port_to_kill ha sido terminado."
done
