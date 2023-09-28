
package main

import (
	"encoding/gob"
	"log"
	"fmt"
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

	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	var request com.Request
	err := decoder.Decode(&request)
	com.CheckError(err)

	if request.Id == -1 {
		fmt.Println("Last client.")
		return
	}

	// Process the request
	reply := com.Reply{
		Id:     request.Id,
		Primes: findPrimes(request.Interval),
	}

	// Send the reply back to the client
	err = encoder.Encode(reply)
	com.CheckError(err)
}

func main() {
	CONN_TYPE := "tcp"
	endpoint := ":30000"

	listener, err := net.Listen(CONN_TYPE, endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	log.Println("***** Listening for new connection in endpoint ", endpoint)

	for {
		conn, err := listener.Accept()
		com.CheckError(err)

		go handleRequest(conn) // Crea una goroutine para manejar cada petición
	}
}