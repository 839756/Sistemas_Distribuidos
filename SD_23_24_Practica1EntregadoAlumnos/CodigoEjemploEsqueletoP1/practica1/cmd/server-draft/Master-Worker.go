package main

import (
	"encoding/gob"
	"log"
	"fmt"
	"net"
	"practica1/com"
	"os/exec"
)

const CONN_TYPE = "tcp"

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

func startWorker(address string) {
	// Hacer que se ejecute un script que compile y que ejecute los workers
	// cmd := exec.Command("ssh","","/home/a839756/SD/practica1/cmd/worker")
}

func handleWorker(address string, client <-chan net.Conn) {
	// Set up the worker address
	endpoint := fmt.Sprintf("%s:29120", address)

	for conn := range client {
		defer conn.Close()

		// New decoder and encoder for client communication
		decoder := gob.NewDecoder(conn)
		encoder := gob.NewEncoder(conn)

		var request com.Request
		err := decoder.Decode(&request)
		com.CheckError(err)

		if request.Id == -1 {
			fmt.Println("Last client.")
		} else {
			// Connect to the worker
			wConn, err := net.Dial("tcp", endpoint)
			com.CheckError(err)

			defer wConn.Close()

			// New decoder and encoder for worker communication
			wDecoder := gob.NewDecoder(wConn)
			wEncoder := gob.NewEncoder(wConn)

			// Send the worker work
			err = wEncoder.Encode(&request)
			com.CheckError(err)

			// Get the worker reply
			var reply com.Request
			err := wDecoder.Decode(&reply)
			com.CheckError(err)

			// Send the reply to the client
			err = encoder.Encode(&reply)
			com.CheckError(err)
		}
	}
}

// Grupos Puertoinicial Puertofinal Máquinainicial Máquinafinal
// 2.5    29120         29129       9              12

func main() {
	listener, err := net.Listen(CONN_TYPE, endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	log.Println("***** Listening for new connection in endpoint ", endpoint)

	client := make(chan net.Conn)

	for i := 10; i <=12; i++ {
		// Setting up IP
		address := fmt.Sprintf("%s%d","192.168.3.",i)
		
		// Start the worker
		go startWorker(address)

		time.Sleep(1 * time.Second)

		// Let the worker work
		go handleWorker(address,client)
	}

	for {
        conn, err := listener.Accept()
        com.CheckError(err)

        client <- conn
    }
}