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
	"log"
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
	log.Printf("Envío para que copien %s\n", mensaje)
	for i := 1; i <= ra.LE; i++ {
		if i != me {

			message.Send(i, Text{mensaje, me})
		}
	}
}

// Responde al que ha enviado el texto que lo ha copiado de forma satisfactoria
func SendReplyToTxt(message *ms.MessageSystem, to int) {
	log.Printf("Hay que mandar una respueta a %d\n", to)
	message.Send(to, TextReply{})
	log.Printf("Se ha mandado la respuesta a %d\n", to)
}

// Espera a que los demás nodos hayan copiado.
func WaitForReply(chtxt chan bool, numEsperar int) {
	for i := numEsperar; i > 0; i-- {
		<-chtxt
	}
}

func Receptor(msg *ms.MessageSystem, chreq chan ra.Request, chrep chan ra.Reply, chCheck chan bool, chtxt chan bool, file *gestorF.Fich) {
	for {
		mensaje := msg.Receive()
		switch tipo := mensaje.(type) {
		case ra.Request:
			chreq <- tipo
		case ra.Reply:
			if tipo.Post {
				log.Printf("El proceso %d ha recibido permiso postergado\n", tipo.Recibido)
			} else {
				log.Printf("El proceso %d ha recibido permiso \n", tipo.Recibido)
			}
			chrep <- tipo
		case CheckPoint:
			chCheck <- true
		case Text:
			file.EscribirFichero(tipo.Text)
			SendReplyToTxt(msg, tipo.Pid)
		case TextReply:
			chtxt <- true
		}
	}
}
