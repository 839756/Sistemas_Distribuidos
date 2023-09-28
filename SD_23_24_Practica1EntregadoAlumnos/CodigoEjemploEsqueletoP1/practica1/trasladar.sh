#!/bin/bash

# Dirección IP base
base_ip="192.168.3."

# Usuario para la conexión SSH
usuario="a839756"

# Carpeta que deseas copiar
carpeta_origen="../practica1"

# Bucle para iterar a través de las direcciones IP
for i in {9..12}
do
  # Dirección IP completa
  ip="${base_ip}${i}"
  
  # Comando SCP para copiar la carpeta practica1 a la máquina remota
  scp -r "$carpeta_origen" "$usuario@$ip:~/"
  
  # Verificar el resultado
  if [ $? -eq 0 ]; then
    echo "Carpeta copiada a $ip correctamente"
  else
    echo "Error al copiar la carpeta a $ip"
  fi
done