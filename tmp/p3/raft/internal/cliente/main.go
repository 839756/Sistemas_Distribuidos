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

	cliente, err := rpc.Dial("tcp", "192.168.3.10:29120")
	check.CheckError(err, "No se ha podido conectar al servidor")

	err = cliente.Call("NodoRaft.SometerOperacionRaft", raft.TipoOperacion{Operacion: "leer", Clave: "0", Valor: "0"}, &reply)
	check.CheckError(err, "Ha fallado la llamada")

	err = cliente.Call("NodoRaft.SometerOperacionRaft", raft.TipoOperacion{Operacion: "escribir", Clave: "0", Valor: "0"}, &reply)
	check.CheckError(err, "Ha fallado la llamada")

	err = cliente.Call("NodoRaft.SometerOperacionRaft", raft.TipoOperacion{Operacion: "leer", Clave: "0", Valor: "0"}, &reply)
	check.CheckError(err, "Ha fallado la llamada")

	fmt.Println("Ha terminando de enviar cosas")
}
