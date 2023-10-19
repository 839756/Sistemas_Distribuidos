#!/bin/bash

    # Iniciar el programa barrera.go en una nueva ventana de terminal usando gnome-terminal
    gnome-terminal -- bash -c "go run barrera.go; read -p 'Presione [Enter] para cerrar esta terminal...'"

    # Darle tiempo al programa barrera.go para que se inicie correctamente
    sleep 4

    # Iniciar el programa lector en una nueva ventana de terminal y navegar a su carpeta primero
    gnome-terminal --working-directory=$(pwd)/cmd/lector -- bash -c "go run main.go 1; read -p 'Presione [Enter] para cerrar esta terminal...'"

    # Iniciar el programa lector en una nueva ventana de terminal y navegar a su carpeta primero
    gnome-terminal --working-directory=$(pwd)/cmd/lector -- bash -c "go run main.go 2; read -p 'Presione [Enter] para cerrar esta terminal...'"

    # Darle tiempo al programa lector para que se inicie correctamente
    sleep 1

    # Iniciar el programa escritor en una nueva ventana de terminal y navegar a su carpeta primero
    gnome-terminal --working-directory=$(pwd)/cmd/escritor -- bash -c "go run main.go 3; read -p 'Presione [Enter] para cerrar esta terminal...'"

    # Iniciar el programa escritor en una nueva ventana de terminal y navegar a su carpeta primero
    gnome-terminal --working-directory=$(pwd)/cmd/escritor -- bash -c "go run main.go 4; read -p 'Presione [Enter] para cerrar esta terminal...'"

    # Darle tiempo al programa escritor para que termine de ejecutar
    sleep 4

    # Comparar el contenido de fichero_1.txt y fichero_2.txt
    if cmp -s "cmd/lector/fichero_1.txt" "cmd/escritor/fichero_3.txt" && cmp -s "cmd/lector/fichero_2.txt" "cmd/escritor/fichero_4.txt"; then
        echo "Los archivos son iguales."
    else
        echo "Los archivos son diferentes."
    fi

    echo "FIN."
    echo "--------------------------"

