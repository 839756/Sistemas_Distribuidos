package main

import (
	"fmt"
	"net/rpc"

	//"log"
	//"crypto/rand"

	"raft/internal/comun/check"
	"raft/internal/raft"
)

func main() {

	var reply raft.ResultadoRemoto

	cliente, err := rpc.Dial("tcp", "localhost:29001")
	check.CheckError(err, "No se ha podido conectar al servidor")

	err = cliente.Call("NodoRaft.SometerOperacionRaft", raft.TipoOperacion{Operacion: "leer", Clave: "1", Valor: "1"}, &reply)
	check.CheckError(err, "Ha fallado la llamada")

	err = cliente.Call("NodoRaft.SometerOperacionRaft", raft.TipoOperacion{Operacion: "escribir", Clave: "2", Valor: "2"}, &reply)
	check.CheckError(err, "Ha fallado la llamada")

	err = cliente.Call("NodoRaft.SometerOperacionRaft", raft.TipoOperacion{Operacion: "leer", Clave: "3", Valor: "3"}, &reply)
	check.CheckError(err, "Ha fallado la llamada")

	fmt.Println("Ha terminando de enviar cosas")
}
