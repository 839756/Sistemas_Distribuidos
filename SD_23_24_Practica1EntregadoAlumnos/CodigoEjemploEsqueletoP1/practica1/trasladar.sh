#!/bin/bash

# Definir la carpeta que deseas copiar
carpeta_a_copiar="practica1"

# Definir el nombre de usuario en el servidor remoto
usuario="a839756"

# Definir la ruta del archivo Go a compilar
archivo_go="practica1/cmd/worker/main.go"

# Iterar sobre las direcciones IP 192.168.3.10, 192.168.3.11 y 192.168.3.12
for ip in {10..12}; do
    servidor="192.168.3.$ip"
    
    # Usar scp para copiar la carpeta al servidor remoto
    scp -r "$carpeta_a_copiar" "$usuario@$servidor:~/"
    
    # Verificar el resultado de la copia
    if [ $? -eq 0 ]; then
        echo "Carpeta copiada a $servidor"
        
        # Compilar el archivo Go en el servidor
        ssh "$usuario@$servidor" "go build $archivo_go"
        
        # Verificar el resultado de la compilación
        if [ $? -eq 0 ]; then
            echo "Archivo Go compilado en $servidor"
        else
            echo "Error al compilar el archivo Go en $servidor"
        fi
    else
        echo "Error al copiar la carpeta a $servidor"
    fi
done
