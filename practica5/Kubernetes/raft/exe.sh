#!/bin/bash

# Abre la primera terminal y ejecuta el comando
gnome-terminal  --title="Terminal 1" -- bash -c "go run cmd/srvraft/main.go 0 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003; read -p 'Presiona Enter para cerrar esta terminal'"

sleep 1

# Abre la segunda terminal y ejecuta el comando
gnome-terminal  --title="Terminal 2" -- bash -c "go run cmd/srvraft/main.go 1 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003; read -p 'Presiona Enter para cerrar esta terminal'"

# Abre la tercera terminal y ejecuta el comando
gnome-terminal  --title="Terminal 3" -- bash -c "go run cmd/srvraft/main.go 2 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003; read -p 'Presiona Enter para cerrar esta terminal'"

sleep 1

gnome-terminal  --title="Terminal Cliente" -- bash -c "go run internal/cliente/main.go"

sleep 1

pid=$(lsof -ti tcp:29003)
  
  if [ -n "$pid" ]; then
    echo "El puerto 29003 está ocupado por el proceso $pid. Matando el proceso..."
    kill -9 $pid
    if [ $? -eq 0 ]; then
      echo "Proceso $pid matado con éxito."
    else
      echo "No se pudo matar el proceso $pid."
    fi
  else
    echo "El puerto $port está libre."
  fi

sleep 2

gnome-terminal  --title="Terminal 3 nueva" -- bash -c "go run cmd/srvraft/main.go 2 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003; read -p 'Presiona Enter para cerrar esta terminal'"

gnome-terminal  --title="Terminal Cliente" -- bash -c "go run internal/cliente/main.go"

sleep 1

./liberar.sh