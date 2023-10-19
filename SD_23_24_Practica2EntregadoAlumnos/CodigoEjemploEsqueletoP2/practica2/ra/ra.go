/*
* AUTOR: Rafael Tolosana Calasanz
* AUTOR: Jorge Leris Lacort - 845647
* AUTOR: Andrei Gabriel Vlasceanu - 839756
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
 */
package ra

import (
	"log"
	"practica2/ms"
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec/vclock"
)

const (
	LE = 4
)

type Request struct {
	Clock     vclock.VClock
	Pid       int
	Operation string
}

type Reply struct{}

type RASharedDB struct {
	OurSeqNum vclock.VClock     // Numero de secuencia enviado del propio nodo
	HigSeqNum vclock.VClock     // El número de secuencia más alto recibido
	OutRepCnt int               // Número de respuestas esperado
	ReqCS     bool              // ¿Está haciendo una peticion?
	RepDefd   []bool            // Nodos los cuales han sido postergados
	ms        *ms.MessageSystem // Tipo mensaje
	done      chan bool         // Canal para confirmar que ha terminado
	chrep     chan bool
	Mutex     sync.Mutex // mutex para proteger concurrencia sobre las variables
	// TODO: completar
	exclude   map[string]map[string]bool // Matriz de exclusión formada por un mapa de mapas de booleanos
	operation string                     // Operación que hace el nodo
	repl      chan Reply
	reqt      chan Request
}

func New(msgs *ms.MessageSystem, me int, usersFile string, operation string, resp chan Reply, pet chan Request) *RASharedDB {

	ra := RASharedDB{vclock.New(), vclock.New(), 0, false, make([]bool, LE), msgs, make(chan bool), make(chan bool),
		sync.Mutex{}, make(map[string]map[string]bool), operation, resp, pet}

	ra.exclude["read"] = make(map[string]bool)
	ra.exclude["read"]["read"] = false
	ra.exclude["read"]["write"] = true

	ra.exclude["write"] = make(map[string]bool)
	ra.exclude["write"]["read"] = true
	ra.exclude["write"]["write"] = true

	// Iniciamos el HigSeqNum a 0
	ra.HigSeqNum.Set(strconv.Itoa(me), 0)

	// Iniciamos go routines para escuchar
	go ra.permission()
	go ra.request()

	return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol() {

	me := ra.ms.WhoSends()
	ra.Mutex.Lock()

	ra.ReqCS = true                                                 // Pide la sección crítica
	ra.OurSeqNum.Set(strconv.Itoa(me), ra.HigSeqNum.LastUpdate()+1) // Actualizamos el reloj

	ra.Mutex.Unlock()

	ra.OutRepCnt = LE - 1 // Numero de respuestas que se esperan

	log.Printf("Mi reloj al enviar el evento es este: ")
	ra.OurSeqNum.PrintVC()
	// Mandamos solicitud a los demás nodos
	for pid := 0; pid < LE; pid++ {
		if pid != me {
			ra.Mutex.Lock()
			ra.ms.Send(pid+1, Request{ra.OurSeqNum, me, ra.operation})
			ra.Mutex.Unlock()
		}
	}

	// Esperamos respuesta de los demás nodos
	<-ra.chrep

}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol() {

	ra.ReqCS = false

	for pid := 0; pid < LE; pid++ { // Se recorren todos los procesos lector/escritor

		if ra.RepDefd[pid] { // for each j ∈ perm_delayedi

			ra.RepDefd[pid] = false
			ra.Mutex.Lock()
			ra.ms.Send(pid+1, Reply{}) // Envia un reply al nodo pid+1
			ra.Mutex.Unlock()
		}
	} //end for

}

// La función maneja la recepción de un mensaje REQUEST (k, j).
// Actualiza el reloj local, verifica si se puede otorgar permiso y responde en consecuencia.
func (ra *RASharedDB) request() {
	for {
		request := <-ra.reqt
		var defer_it bool

		log.Printf("El reloj recibido es el siguiente: ")
		request.Clock.PrintVC()

		log.Printf("Mi reloj es el siguiente: ")
		ra.HigSeqNum.PrintVC()

		ra.HigSeqNum.Merge(request.Clock) // Actualizamos relojes con el que hemos recibido

		log.Printf("El reloj combiando es el siguiente: ")
		ra.HigSeqNum.PrintVC()

		me := ra.ms.WhoSends()
		him := request.Pid

		myClock, myErr := ra.HigSeqNum.FindTicks(strconv.Itoa(me))
		if !myErr {
			log.Println("Error en la lectura del propio reloj")
		}
		hisClock, hisErr := ra.HigSeqNum.FindTicks(strconv.Itoa(him))
		if !hisErr {
			log.Println("Error en la lectura del reloj que manda el nodo")
		}
		ra.Mutex.Lock() // Consulta de variables compartidas
		defer_it = ra.ReqCS && decidePriority(myClock, hisClock, me, him) && ra.exclude[ra.operation][request.Operation]
		ra.Mutex.Unlock()

		if defer_it {
			ra.RepDefd[him-1] = true
		} else {
			ra.ms.Send(him, Reply{})
		}
	}
}

func decidePriority(myClock uint64, hisClock uint64, me int, him int) bool {
	if myClock > hisClock {
		return true
	} else if myClock == hisClock {
		return me < him
	} else {
		return false
	}

}

// La función maneja la recepción de un mensaje PERMISSION(j).
// Elimina j de la lista de esperas (waiting_fromi).
func (ra *RASharedDB) permission() {
	for {
		<-ra.repl
		ra.OutRepCnt--
		if ra.OutRepCnt == 0 {
			ra.chrep <- true
		}
	}

}

func (ra *RASharedDB) Stop() {
	ra.ms.Stop()
	ra.done <- true
}
