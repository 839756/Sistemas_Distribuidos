/*
* AUTOR: Jorge Leris Lacort - 845647
* AUTOR: Andrei Gabriel Vlasceanu - 839756
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: octubre de 2023
* FICHERO: main.go
* DESCRIPCIÓN: Implementación del escritor del sistema
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
	"time"
)

func escritor(file *gestorF.Fich, ricart *ra.RASharedDB, message *ms.MessageSystem, me int, wait *sync.WaitGroup, chRpl chan bool) {
	defer wait.Done()

	for j := 0; j < 100; j++ {
		ricart.PreProtocol()
		// Escribimos en el fichero

		file.EscribirFichero(strconv.Itoa(j) + " ")

		log.Printf("Contenido escrito por el proceso %d: %d\n", me, j)
		// Enviamos un mensaje para que se actualicen los ficheros de los demás procesos
		receptor.SendText(message, j, me)
		// Espera a recibir todos los
		// receptor.WaitForReply(chRpl, ra.LE-1)
		for i := ra.LE - 1; i > 0; i-- {
			<-chRpl
		}
		ricart.PostProtocol()
	}
	// Crea un temporizador para esperar segundos
	duration := 5 * time.Second
	timer := time.NewTimer(duration)

	// Espera los  segundos para terminar la ejecucion y que se actualice todo
	<-timer.C
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	myPid := os.Args[1]
	me, _ := strconv.Atoi(myPid)
	log.Printf("Escritor con pid %d en marcha\n", me)
	// Se crea la copia del fichero
	fichero := "fichero_" + myPid + ".txt"
	file := gestorF.CrearFichero(fichero)

	usersFile := "../../ms/users.txt" // Fichero con dirección de las demás máquinas

	tipoDeMensajes := []ms.Message{ra.Request{}, ra.Reply{}, receptor.CheckPoint{}, receptor.Text{}}

	message := ms.New(me, usersFile, tipoDeMensajes)
	// Creamos los canales para comunicarse con el algoritmo RA
	chReq := make(chan ra.Request)
	chRep := make(chan ra.Reply)
	chCheck := make(chan bool)
	chTxtRepl := make(chan bool)
	// Iniciamos el receptor de mensaje
	ricart := ra.New(&message, me, usersFile, "write", chRep, chReq)

	go receptor.Receptor(&message, chReq, chRep, chCheck, chTxtRepl, file)
	log.Println("Receptor iniciado")

	message.Send(ra.LE+1, receptor.CheckPoint{})
	<-chCheck
	var wait sync.WaitGroup
	wait.Add(1)
	go escritor(file, ricart, &message, me, &wait, chTxtRepl)
	wait.Wait()
}
