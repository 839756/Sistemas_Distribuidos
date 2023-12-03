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

	fmt.Println("Comienza el cliente")
	fmt.Println()

	time.Sleep(10000 * time.Millisecond) //Espera a que se pongan los servidores en marcha

	lider := obtenerLider(3, nodo) //Se obtiene el lider
	if lider == -1 {
		fmt.Println("No se ha elegido lider")
	}else{
		fmt.Println("El lider es ", lider)
	}

	//Se somete una operacion
	err := nodo[lider].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{Operacion: "leer", Clave: "1", Valor: ""}, &reply, 1000*time.Millisecond)
	check.CheckError(err, "Ha fallado la llamada")

	fmt.Println("Operacion basica sometida")
	fmt.Println("El indice resultante es:", reply.IndiceRegistro, "y el mandato resultante es: ", reply.Mandato)
	fmt.Println()

	//Se somete una operacion tras una caida de seguidor
	for i := 0; i < 3; i++ {
		if i != lider {
			pararNodo(i, nodo)
			fmt.Println("Se ha parado el nodo: ", i)

			err = nodo[lider].CallTimeout("NodoRaft.SometerOperacionRaft",
				raft.TipoOperacion{Operacion: "escribir", Clave: "2", Valor: "2"}, &reply, 1000*time.Millisecond)
			check.CheckError(err, "Ha fallado la llamada")

			fmt.Println("Operacion con caida de seguidor sometida")
			fmt.Println("El indice resultante es:", reply.IndiceRegistro, "y el mandato resultante es: ", reply.Mandato)

			break
		}
	}
	fmt.Println()

	//Se somete una operacion tras una caida del lider
	pararNodo(lider, nodo)
	fmt.Println("Se ha parado el actual lider: ", lider)

	lider = obtenerLider(3, nodo)
	if lider == -1 {
		fmt.Println("No se ha elegido lider tras la caida del anterior")
	}else{
		fmt.Println("Ahora el lider es ", lider)
	}

	err = nodo[lider].CallTimeout("NodoRaft.SometerOperacionRaft",
		raft.TipoOperacion{Operacion: "leer", Clave: "3", Valor: ""}, &reply, 1000*time.Millisecond)
	check.CheckError(err, "Ha fallado la llamada")

	fmt.Println("Operacion con caida de lider sometida")
	fmt.Println("El indice resultante es:", reply.IndiceRegistro, "y el mandato resultante es: ", reply.Mandato)

	fmt.Println()
	fmt.Println("Ha terminando de enviar cosas")
}

// --------------------------------------------------------------------------
// FUNCIONES DE APOYO
// Funcion que obtiene el lider
func obtenerLider(numreplicas int, nodo []rpctimeout.HostPort) int {
	var reply raft.EstadoRemoto

	time.Sleep(3000 * time.Millisecond)
	for i := 0; i < numreplicas; i++ {
		err := nodo[i].CallTimeout("NodoRaft.ObtenerEstadoNodo",
			raft.Vacio{}, &reply, 100*time.Millisecond)
		check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")
		if reply.EsLider {
			return reply.IdLider
		}
	}
	return -1
}

//Funcion para parar un nodo
func pararNodo(replica int, nodo []rpctimeout.HostPort) {
	var reply raft.EstadoRemoto

	err := nodo[replica].CallTimeout("NodoRaft.ParaNodo",
		raft.Vacio{}, &reply, 100*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ParaNodo")

	time.Sleep(5000 * time.Millisecond)
}
