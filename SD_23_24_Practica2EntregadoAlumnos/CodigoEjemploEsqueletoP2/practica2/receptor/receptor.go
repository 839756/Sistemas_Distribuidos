/*
* AUTOR: Jorge Leris Lacort - 845647
* AUTOR: Andrei Gabriel Vlasceanu - 839756
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: octubre de 2023
* FICHERO: receptor.go
* DESCRIPCIÓN: Implementación de la función encargada de recibir los mensajes del sistema
 */

package receptor

import (
	"practica2/gestorF"
	"practica2/ms"
	"practica2/ra"
)

type CheckPoint struct{}

type Text struct {
	Text string
}

func Receptor(msg *ms.MessageSystem, chreq chan ra.Request, chrep chan ra.Reply, chCheck chan bool, file string) {
	for {
		mensaje := msg.Receive()
		switch tipo := mensaje.(type) {
		case ra.Request:
			chreq <- tipo
		case ra.Reply:
			chrep <- tipo
		case CheckPoint:
			chCheck <- true
		case Text:
			gestorF.EscribirFichero(file, tipo.Text)
		}
	}
}
