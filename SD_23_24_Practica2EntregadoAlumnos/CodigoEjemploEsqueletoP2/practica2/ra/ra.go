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
    "practica2/ms"
    "sync"
    "github.com/DistributedClocks/GoVector/govec/vclock"
)

const(
    LE = 4
)

type Request struct{
    Clock   vclock.VClock
    Pid     int

}

type Reply struct{}

type RASharedDB struct {
    OurSeqNum   vclock.VClock		// Numero de secuencia enviado del propio nodo
    HigSeqNum   int				    // El número de secuencia más alto recibido
    OutRepCnt   int				    // Número de respuestas esperado
    ReqCS       bool				// ¿Está haciendo una peticion?
    RepDefd     []bool		    	// Nodos los cuales han sido postergados
    ms          *MessageSystem		// Tipo mensaje
    done        chan bool			// Canal para confirmar que ha terminado
    chrep       chan bool			
    Mutex       sync.Mutex // mutex para proteger concurrencia sobre las variables
    // TODO: completar
    exclude    map[string]map[string] bool // Matriz de exclusión formada por un mapa de mapas de booleanos
}

func New(me int, usersFile string) (*RASharedDB) {
    messageTypes := []Message{Request, Reply}
    msgs = ms.New(me, usersFile string, messageTypes)
    ra := RASharedDB{0, 0, 0, false, []int{}, &msgs,  make(chan bool),  make(chan bool), &sync.Mutex{}, make(map[string]map[string]bool)}
    // TODO completar
    ra.exclude["read"]["read"] = false
    ra.exclude["read"]["write"] = true
    ra.exclude["write"]["read"] = true
    ra.exclude["write"]["write"] = true
    
    return &ra
}

//Pre: Verdad
//Post: Realiza  el  PreProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol(){
    // TODO completar
}

//Pre: Verdad
//Post: Realiza  el  PostProtocol  para el  algoritmo de
//      Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol(){

    ra.ReqCS = false                    // cs_statei ← out;
    
    for pid := 0; pid <= LE; pid++ {    // Se recorren todos los procesos lector/escritor
        
        if ra.RepDefd[pid] {            // for each j ∈ perm_delayedi
            
            ra.Mutex.Lock()
            ra.ms.Send(pid+1, Reply{})  // do send PERMISSION(i) to pj 
            ra.Mutex.Unlock()
            ra.RepDefd[pid] = false;    // perm_delayedi ← ∅
        }
    }  //end for

}

// La función maneja la recepción de un mensaje REQUEST (k, j).
// Actualiza el reloj local, verifica si se puede otorgar permiso y responde en consecuencia.
func (ra *RASharedDB) request(){
    // TODO completar
}

// La función maneja la recepción de un mensaje PERMISSION(j).
// Elimina j de la lista de esperas (waiting_fromi).
func (ra *RASharedDB) permission(){
    // TODO completar
}

func (ra *RASharedDB) Stop(){
    ra.ms.Stop()
    ra.done <- true
}
