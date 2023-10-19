/*
* AUTOR: Jorge Leris Lacort - 845647
* AUTOR: Andrei Gabriel Vlasceanu - 839756
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: octubre de 2023
* FICHERO: main.go
* DESCRIPCIÓN: Implementación del lector del sistema
 */

package main

import (
	"log"
	"os"
	"practica2/gestorF"
	"practica2/ms"
	"practica2/ra"
	"practica2/receptor"
	"strconv"
	"sync"
)

func lector(fichero string, ricart *ra.RASharedDB, wait *sync.WaitGroup) {
	defer wait.Done()

	for i := 0; i < 10; i++ {
		ricart.PreProtocol()
		//Leer en el fichero
		gestorF.LeerFichero(fichero)

		ricart.PostProtocol()
	}
	/*
		// Crea un temporizador para esperar 5 segundos
		duration := 5 * time.Second
		timer := time.NewTimer(duration)

		// Espera los 5 segundos para terminar la ejecucion y que se actualice todo
		<-timer.C*/
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds) // Set log config

	myPid := os.Args[1]
	me, _ := strconv.Atoi(myPid)
	log.Printf("Lector con pid %d en marcha\n", me)
	// Se crea la copia del fichero
	fichero := "fichero_" + myPid + ".txt"
	gestorF.CrearFichero(fichero)

	usersFile := "../../ms/users.txt" // Fichero con dirección de las demás máquinas

	tipoDeMensajes := []ms.Message{ra.Request{}, ra.Reply{}, receptor.CheckPoint{}, receptor.Text{}, receptor.TextReply{}}

	message := ms.New(me, usersFile, tipoDeMensajes)
	// Creamos los canales para comunicarse con el algoritmo RA
	chReq := make(chan ra.Request)
	chRep := make(chan ra.Reply)
	chCheck := make(chan bool)
	chtext := make(chan bool)
	// Iniciamos el receptor de mensaje
	go receptor.Receptor(&message, chReq, chRep, chCheck, chtext, fichero)
	log.Println("Receptor iniciado")

	ricart := ra.New(&message, me, usersFile, "read", chRep, chReq)

	message.Send(ra.LE+1, receptor.CheckPoint{})

	log.Println("Esperando barrera")
	<-chCheck
	log.Println("Barrera pasada")

	var wait sync.WaitGroup
	wait.Add(1)
	go lector(fichero, ricart, &wait) //PONER FEEDBACK DE LO LEIDO
	wait.Wait()
	// Terminar cuando los demás procesos terminen también
	message.Send(ra.LE+1, receptor.CheckPoint{})
	<-chCheck
}
