package main

import (
	"fmt"
	"strconv"
	"time"

	//"log"
	//"crypto/rand"

	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
)

func main() {
	var reply raft.ResultadoRemoto
	var dir_servidores []string
	var nodo []rpctimeout.HostPort

	ss_service := ".ss-service.default.svc.cluster.local:6000"

	for i := 0; i < 3; i++ {
		dir := "ss-" + strconv.Itoa(i) + ss_service
		dir_servidores = append(dir_servidores, dir)
	}
	for _, endPoint := range dir_servidores {
		nodo = append(nodo, rpctimeout.HostPort(endPoint))
	}

	time.Sleep(4000 * time.Millisecond) //Espera a que se pongan los servidores en marcha

	err := nodo[0].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{Operacion: "leer", Clave: "1", Valor: ""}, &reply, 2000*time.Millisecond)
	check.CheckError(err, "Ha fallado la llamada")

	err = nodo[1].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{Operacion: "escribir", Clave: "2", Valor: "2"}, &reply, 2000*time.Millisecond)
	check.CheckError(err, "Ha fallado la llamada")

	err = nodo[2].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{Operacion: "leer", Clave: "3", Valor: ""}, &reply, 2000*time.Millisecond)
	check.CheckError(err, "Ha fallado la llamada")

	fmt.Println("Ha terminando de enviar cosas")
}
