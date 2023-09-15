
 package main

 import (
	 //"encoding/gob"
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
 
 // COMPLETAR EL SERVIDOR  .....
 func main() {
 
	 CONN_TYPE := "tcp"
	 endpoint := ":30000"
 	 
		 for {
			listener, err := net.Listen(CONN_TYPE, endpoint)
	 		com.CheckError(err)
 
	 		log.SetFlags(log.Lshortfile | log.Lmicroseconds)

		 	log.Println("***** Listening for new connection in endpoint ", endpoint)
		 	conn, err := listener.Accept()
		 	defer conn.Close()
		 	com.CheckError(err)
			}	 
 }