package main

import (
	"flag"
	"log"
	"net"
	"os"
)

func main() {
	var tcpAddr string
	flag.StringVar(&tcpAddr, "tcp", "", "Run as a TCP client to connect target address.")
	flag.Parse()
	if tcpAddr == "" {
		log.Fatalf("You must fill the TCP target address to connect.")
	}

	// Connect TCP target address
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
}
