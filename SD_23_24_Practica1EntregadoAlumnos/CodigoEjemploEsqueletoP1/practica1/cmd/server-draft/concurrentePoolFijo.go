package main

import (
    "net"
    "practica1/com"
    "practica1/utils"
)

func poolHandle(client <-chan net.Conn){
	for conn := range client{
		handleRequest(conn)
	}
}

func main() {
	CONN_TYPE := "tcp"
	endpoint := ":30000"

	listener, err := net.Listen(CONN_TYPE, endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	log.Println("***** Listening for new connection in endpoint ", endpoint)

	client := make(chan net.Conn)

	poolSize := 5
    
    for i := 0; i < poolSize; i++ {
        go poolHandle(client)
    }

	for {
        conn, err := listener.Accept()
        com.CheckError(err)

        client <- conn
    }


}