#!/bin/bash

# Ruta completa a la carpeta raft
RAFT_DIR=~/tmp/p4/raft

# Comprobar si la carpeta existe y eliminarla
if [ -d "$RAFT_DIR" ]; then
    rm -rf "$RAFT_DIR"
    echo "La carpeta 'raft' ha sido eliminada."
else
    echo "La carpeta 'raft' no existe en '$RAFT_DIR'."
fi

cp -r ../. ~/tmp/p4/

# Comprobar si el comando de copia se ejecutó correctamente
if [ $? -eq 0 ]; then
    echo "La copia se realizó correctamente."
else
    echo "Hubo un error al realizar la copia."
fi
