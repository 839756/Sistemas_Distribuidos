//Autor: Jorge Leris Lacort - 845647
//Autor: Andrei Gabriel Vlasceanu - 839756

package main

import (
	"encoding/gob"
	"fmt"
	"log"
    "net"
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

func poolHandle(client <-chan net.Conn){
	// Loop through each connection received from the 'client' channel.
	for conn := range client{
		// Call the 'handleRequest' function to process the connection.
		handleRequest(conn)
	}
}

func main() {
	CONN_TYPE := "tcp"
	endpoint := ":29120"

	listener, err := net.Listen(CONN_TYPE, endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	log.Println("***** Listening for new connection in endpoint ", endpoint)

	client := make(chan net.Conn)

	poolSize := 5
    
    for i := 0; i < poolSize; i++ {
		// Create a pool of goroutines to handle incoming connections.
        go poolHandle(client)
    }

	for {
        conn, err := listener.Accept()
        com.CheckError(err)
		log.Println("Client accepted")
    	// Send the newly accepted connection to one of the goroutines in the pool.
        client <- conn
    }


}