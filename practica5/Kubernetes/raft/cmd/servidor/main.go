package main

import (
	//"errors"
	"fmt"
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

	// obtener entero de indice de este nodo
	me, err := strconv.Atoi(os.Args[1])
	check.CheckError(err, "Main, mal numero entero de indice de nodo:")

	datos := make(map[string]string)
	ss_service := ".ss-service.default.svc.cluster.local:6000"
	var dir_servidores []string
	var nodos []rpctimeout.HostPort

	for i := 0; i < 3; i++ {
		dir := "ss-" + strconv.Itoa(i) + ss_service
		dir_servidores = append(dir_servidores, dir)
	}

	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range dir_servidores {
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	AplicaOpChan := make(chan raft.AplicaOperacion, 1000)

	// Parte Servidor
	nr := raft.NuevoNodo(nodos, me, AplicaOpChan)
	rpc.Register(nr)

	go aplicaOp(datos, AplicaOpChan)

	fmt.Println("Replica escucha en :", me, " de ", os.Args[2:])

	l, err := net.Listen("tcp", os.Args[2:][me])
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
