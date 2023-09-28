//Autor: Jorge Leris Lacort - 845647
//Autor: Andrei Gabriel Vlasceanu - 839756


package main

import (
	"encoding/gob"
	"fmt"
    "net"
	"log"
    "practica1/com"
)

// PRE: verdad = !foundDivisor
// POST: IsPrime devuelve verdad si n es primo y falso en caso contrario
func isPrime(n int) (foundDivisor bool) {
	foundDivisor = false
	for i := 2; (i < n) && !foundDivisor; i++ {
		foundDivisor = (n%i == 0)
	}
	return !foundDivisor
}

// PRE: interval.A < interval.B
// POST: FindPrimes devuelve todos los números primos comprendidos en el
//
//	intervalo [interval.A, interval.B]
func findPrimes(interval com.TPInterval) (primes []int) {
	for i := interval.Min; i <= interval.Max; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}


func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Create a decoder and encoder for the network connection
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	// Declare a variable to hold the incoming request
	var request com.Request
	// Decode the incoming request and check for errors
	err := decoder.Decode(&request)
	com.CheckError(err)

	// Check if the request ID is -1, indicating it's the last client
	if request.Id == -1 {
		fmt.Println("Last client.")
		return
	}

	// Prepare a reply based on the request, including finding prime numbers
	reply := com.Reply{
		Id:     request.Id,
		Primes: findPrimes(request.Interval),
	}

	// Transmit the response
	err = encoder.Encode(reply)
	com.CheckError(err)
}


func main() {
	CONN_TYPE := "tcp"
	endpoint := ":29120"

	listener, err := net.Listen(CONN_TYPE, endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	log.Println("***** Listening for new connection in endpoint ", endpoint)

	for {
		conn, err := listener.Accept()
		com.CheckError(err)

		go handleRequest(conn) // Create a goroutine to handle each request.
	}
}