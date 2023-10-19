/*
* AUTOR: Jorge Leris Lacort - 845647
* AUTOR: Andrei Gabriel Vlasceanu - 839756
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: octubre de 2023
* FICHERO: gestorF.go
* DESCRIPCIÓN: Implementación de las dos operaciones necesarias para el proceso gestor de ficheros.
 */

package gestorF

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Fich struct {
	nombre string     // Nombre del fichero
	Mutex  sync.Mutex // Mutex
}

//La operacion crearFichero se encarga de crear un fichero vacio
func CrearFichero(nombreArchivo string) *Fich {
	_, err := os.Create(nombreArchivo)
	if err != nil {
		// Si se produce un error al crear el archivo, mostrar un mensaje y salir del programa.
		fmt.Println("Se ha producido un error al crear el fichero")
		os.Exit(1)
	}
	file := Fich{nombreArchivo, sync.Mutex{}}
	return &file
}

//La operacion LeerFichero devuelve el contenido completo del fichero de texto
func (file *Fich) LeerFichero() (string, error) {
	// Lectura en exlusión mutua
	//file.Mutex.Lock()
	//defer file.Mutex.Unlock()
	// Se lee el contenido del archivo especificado por "nombreArchivo".
	contenido, err := ioutil.ReadFile(file.nombre)

	if err != nil {
		// Si se produjo un error, imprimir un mensaje de error y retornar un string vacío y el error.
		fmt.Println("Se ha producido un error en la lectura del fichero")
		return "", err
	}

	// Si no hubo errores, convertir el contenido del archivo en una cadena y retornarla junto con un valor nulo de error.
	return string(contenido), nil
}

//La operacion EscribirFichero añade al final del fichero de texto un fragmento
//Ademas la operacion de escritura se tiene que hacer de tal forma que el proceso escritor escribe en su fichero,
//pero tambien actualiza los ficheros de todos los procesos lectores y escritores, de manera que todas las copias seran iguales.
func (file *Fich) EscribirFichero(fragmento string) {

	// Edición en exclusión mutua
	//file.Mutex.Lock()
	//defer file.Mutex.Unlock()
	// Abrir el archivo en modo escritura, anexando al final o creando si no existe, con persimos
	// de escritura y lectura para el propietario y el grupo.
	archivo, err := os.OpenFile(file.nombre, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		// Si se produce un error al abrir el archivo, mostrar un mensaje y salir del programa.
		fmt.Println("Se ha producido un error al abrir el fichero")
		os.Exit(1)
	}

	// Se asegura de que el archivo se cierre al finalizar la función.
	defer archivo.Close()

	// Escribir el fragmento en el archivo.
	_, err = archivo.WriteString(fragmento)
	if err != nil {
		// Si se produce un error al escribir en el archivo, mostrar un mensaje y salir del programa.
		fmt.Println("Se ha producido un error al escribir en el fichero")
		os.Exit(1)
	}

}
