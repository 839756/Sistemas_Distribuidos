#!/bin/bash

for port in 29001 29002 29003; do
  # Buscar el proceso que está utilizando el puerto
  pid=$(lsof -ti tcp:$port)
  
  if [ -n "$pid" ]; then
    echo "El puerto $port está ocupado por el proceso $pid. Matando el proceso..."
    kill -9 $pid
    if [ $? -eq 0 ]; then
      echo "Proceso $pid matado con éxito."
    else
      echo "No se pudo matar el proceso $pid."
    fi
  else
    echo "El puerto $port está libre."
  fi
done

