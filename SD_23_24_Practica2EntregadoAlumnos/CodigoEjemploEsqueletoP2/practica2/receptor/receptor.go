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
	"strconv"
)

type CheckPoint struct{}

type Text struct {
	Text string
	Pid  int
}

type TextReply struct{}

// Función que actualiza los ficheros de los demás nodos
func SendText(message *ms.MessageSystem, numMessage int, me int) {
	mensaje := strconv.Itoa(numMessage) + " "
	for i := 1; i <= ra.LE; i++ {
		if i != me {
			message.Send(i, Text{mensaje, me})
		}
	}
}

func SendReplyToTxt(message *ms.MessageSystem, to int) {
	message.Send(to, TextReply{})
}

func WaitForReply(chTxtRpl chan bool, numEsperar int) {
	for i := numEsperar; i > 0; i-- {
		<-chTxtRpl
	}
}

func Receptor(msg *ms.MessageSystem, chreq chan ra.Request, chrep chan ra.Reply, chCheck chan bool, chTxtRpl chan bool, file *gestorF.Fich) {
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
			file.EscribirFichero(tipo.Text)
			SendReplyToTxt(msg, tipo.Pid)
		case TextReply:
			chTxtRpl <- true
		}
	}
}
