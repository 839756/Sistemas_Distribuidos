package main

import (
    "net"
    "practica1/com"
    "practica1/utils"
)

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
		
		handleRequest(conn)	 
	}
}