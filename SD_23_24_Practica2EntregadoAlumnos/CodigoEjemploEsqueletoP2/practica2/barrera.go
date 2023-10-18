/*
* AUTOR: Jorge Leris Lacort - 845647
* AUTOR: Andrei Gabriel Vlasceanu - 839756
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: octubre de 2023
* FICHERO: barrera.go
* DESCRIPCIÓN: Implementación de la barrera distribuida del sistema
 */

package main

import (
	"log"
	"practica2/ms"
	"practica2/ra"
	"practica2/receptor"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	usersFile := "./ms/users.txt"
	messages := []ms.Message{receptor.CheckPoint{}}
	msg := ms.New(ra.LE+1, usersFile, messages)
	for i := 1; i <= ra.LE; i++ {
		_ = msg.Receive()
		log.Println("Proceso en la barrera faltan %d\n", ra.LE-i)
	}
	log.Println("Todos los procesos han llegado a la barrera")
	for i := 1; i <= ra.LE; i++ {
		msg.Send(i, receptor.CheckPoint{})
	}
}
