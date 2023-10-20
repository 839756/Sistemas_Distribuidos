#!/bin/bash
USER="a839756"

# Variables para la primera máquina
IPS=("192.168.3.9" "192.168.3.10" "192.168.3.11" "192.168.3.12")
LECT_DIR="~/practica2/cmd/lector"
ESCT_DIR="~/practica2/cmd/escritor"
FILES=("fichero_1.txt" "fichero_2.txt" "fichero_3.txt" "fichero_4.txt") 

# Usamos ssh para acceder a cada servidor y obtener el contenido de los archivos
file_contents_1=$(ssh "$USER@${IPS[0]}" "cat $LECT_DIR/${FILES[0]}")
file_contents_2=$(ssh "$USER@${IPS[1]}" "cat $LECT_DIR/${FILES[1]}")
file_contents_3=$(ssh "$USER@${IPS[2]}" "cat $ESCT_DIR/${FILES[2]}")
file_contents_4=$(ssh "$USER@${IPS[3]}" "cat $EXCT_DIR/${FILES[3]}")

# Comparamos los contenidos de los archivos
if [ "$file_contents_1" = "$file_contents_2" ] && [ "$file_contents_2" = "$file_contents_3" ] && [ "$file_contents_3" = "$file_contents_4" ]; then
  echo "Los contenidos de los archivos son iguales en las cuatro máquinas."
else
  echo "Los contenidos de los archivos son diferentes en al menos una de las máquinas."
fi
