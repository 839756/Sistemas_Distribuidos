#!/bin/bash

# Define las direcciones IP de las máquinas remotas
remote_ips=("192.168.3.9" "192.168.3.10" "192.168.3.11" "192.168.3.12")

# Captura el nombre de usuario actual de la terminal
remote_user="$USER"

# Define el puerto que deseas verificar y matar
port_to_check=29120

# Recorre las direcciones IP y ejecuta el comando para matar los procesos si es necesario
for ip in "${remote_ips[@]}"; do
  # Utiliza SSH para conectarte a la máquina remota y ejecutar el comando lsof
  ssh "$remote_user@$ip" "lsof -i :$port_to_check" > /tmp/processes_to_kill
  # Verifica si se encontraron procesos en el puerto
  if [ -s /tmp/processes_to_kill ]; then
    # Extrae los IDs de proceso de los resultados y los mata
    pids=$(awk '$1 != "COMMAND" {print $2}' /tmp/processes_to_kill)
    ssh "$remote_user@$ip" "kill -9 $pids"
    echo "Procesos en $ip en el puerto $port_to_check han sido terminados."
  else
    echo "No se encontraron procesos en $ip en el puerto $port_to_check."
  fi
done

port_to_check=29121

# Utiliza SSH para conectarte a la máquina remota y ejecutar el comando lsof
  ssh "$remote_user@192.168.3.9" "lsof -i :$port_to_check" > /tmp/processes_to_kill
  # Verifica si se encontraron procesos en el puerto
  if [ -s /tmp/processes_to_kill ]; then
    # Extrae los IDs de proceso de los resultados y los mata
    pids=$(awk '$1 != "COMMAND" {print $2}' /tmp/processes_to_kill)
    ssh "$remote_user@192.168.3.9" "kill -9 $pids"
    echo "Procesos en $ip en el puerto $port_to_check han sido terminados."
  else
    echo "No se encontraron procesos en $ip en el puerto $port_to_check."
  fi