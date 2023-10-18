package receptor

import (
	"practica2/gestorF"
	"practica2/ms"
	"practica2/ra"
)

type CheckPoint struct{}

type Text struct {
	text string
}

func receptor(msg *ms.MessageSystem, chreq chan ra.Request, chrep chan ra.Reply, chCheck chan bool, file string) {
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
			gestorF.EscribirFichero(file, tipo.text)
		}
	}
}
