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
)

func escritor(fichero string, contenido string, ricart *ra.RASharedDB, message *ms.MessageSystem, me int, wait *sync.WaitGroup) {
	defer wait.Done()

	for j := 0; j < 10; j++ {
		ricart.PreProtocol()
		//Escribimos en el fichero
		gestorF.EscribirFichero(fichero, strconv.Itoa(j)+" ")
		//Enviamos un mensaje para que se actualicen los ficheros de los demás procesos
		for i := 1; i <= ra.LE; i++ {
			if i != me {
				message.Send(i, receptor.Text{Text: strconv.Itoa(j) + " "})
			}
		}
		ricart.PostProtocol()
	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	myPid := os.Args[1]
	me, _ := strconv.Atoi(myPid)
	log.Printf("Escritor con pid %d en marcha\n", me)
	// Se crea la copia del fichero
	fichero := "fichero_" + myPid + ".txt"
	gestorF.CrearFichero(fichero)

	usersFile := "../../ms/users.txt" // Fichero con dirección de las demás máquinas

	tipoDeMensajes := []ms.Message{ra.Request{}, ra.Reply{}, receptor.CheckPoint{}, receptor.Text{}}

	message := ms.New(me, usersFile, tipoDeMensajes)
	// Creamos los canales para comunicarse con el algoritmo RA
	chReq := make(chan ra.Request)
	chRep := make(chan ra.Reply)
	chCheck := make(chan bool)
	// Iniciamos el receptor de mensaje
	go receptor.Receptor(&message, chReq, chRep, chCheck, fichero)
	log.Println("Receptor iniciado")

	ricart := ra.New(&message, me, usersFile, "write", chRep, chReq)

	message.Send(ra.LE+1, receptor.CheckPoint{})
	<-chCheck
	var wait sync.WaitGroup
	wait.Add(1)
	go escritor(fichero, "33 ", ricart, &message, me, &wait)
	wait.Wait()
}
