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
	"time"
)

func lector(file *gestorF.Fich, ricart *ra.RASharedDB, wait *sync.WaitGroup) {
	defer wait.Done()

	for j := 0; j < 100; j++ {
		ricart.PreProtocol()
		//Leer en el fichero
		contenido, _ := file.LeerFichero()
		log.Printf("Contenido leido: %s\n", contenido)
		ricart.PostProtocol()
	}
	// Crea un temporizador para esperar 5 segundos
	duration := 5 * time.Second
	timer := time.NewTimer(duration)

	// Espera los 5 segundos para terminar la ejecucion y que se actualice todo
	<-timer.C
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds) // Set log config

	myPid := os.Args[1]
	me, _ := strconv.Atoi(myPid)
	log.Printf("Lector con pid %d en marcha\n", me)
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
	go receptor.Receptor(&message, chReq, chRep, chCheck, chTxtRepl, file)
	log.Println("Receptor iniciado")

	ricart := ra.New(&message, me, usersFile, "read", chRep, chReq)

	message.Send(ra.LE+1, receptor.CheckPoint{})

	log.Println("Esperando barrera")
	<-chCheck
	log.Println("Barrera pasada")

	var wait sync.WaitGroup
	wait.Add(1)
	go lector(file, ricart, &wait)
	wait.Wait()
}
