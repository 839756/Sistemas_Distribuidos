package main

import (
	//"errors"
	"fmt"
	"strings"
	//"log"
	"net"
	"net/rpc"
	"os"
	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
	"strconv"
	//"time"
)

func main() {

	var nodos []rpctimeout.HostPort
	var dir_servidores []string
	dns:= ".ss-service.default.svc.cluster.local:6000"
	// obtener entero de indice de este nodo

	name := strings.Split(os.Args[1], "-")[0]
	me, err := strconv.Atoi(strings.Split(os.Args[1], "-")[1])
	check.CheckError(err, "Main, mal numero entero de indice de nodo:")

	
	for i := 0; i < 3; i++{
		dir := name + "-" + strconv.Itoa(i) + dns
		dir_servidores = append(dir_servidores,dir)
	}
	
	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range dir_servidores {
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	datos := make(map[string]string)
	AplicaOpChan := make(chan raft.AplicaOperacion, 1000)

	// Parte Servidor
	nr := raft.NuevoNodo(nodos, me, AplicaOpChan)
	rpc.Register(nr)

	go aplicaOp(datos, AplicaOpChan)

	fmt.Println("Replica escucha en :", me, " de ", dir_servidores)

	l, err := net.Listen("tcp", dir_servidores[me])
	check.CheckError(err, "Main listen error:")

	rpc.Accept(l)
}

func aplicaOp(datos map[string]string, AplicaOpChan chan raft.AplicaOperacion) {
	for {
		op := <-AplicaOpChan

		if op.Operacion.Operacion == "leer" {
			op.Operacion.Valor = datos[op.Operacion.Valor]
		} else if op.Operacion.Operacion == "escribir" {
			datos[op.Operacion.Valor] = op.Operacion.Valor
		}
	}
}
