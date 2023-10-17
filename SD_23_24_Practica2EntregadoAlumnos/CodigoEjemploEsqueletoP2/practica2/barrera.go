package main

import (
	"practica2/ra"
	"practica2/ms"
	"practica2/receptor"
	"log"
)


func main(){
	log.SetFlags ( log.Lshortfile | log.Lmicroseconds )
	usersFile := "./ms/users.txt"
	messages := []ms.Message{receptor.CheckPoint}
	msg := ms.New(ra.LE + 1,usersFile, messages)
	for i := 1; i <= ra.LE; i++ {
		_ = msg.Receive()
		log.Println("Proceso en la barrera faltan %d\n", ra.LE - i)
	}
	log.Println("Todos los procesos han llegado a la barrera")
	for i := 1; i <= ra.LE; i++ {
		msg.Send(i, receptor.CheckPoint)
	}
}