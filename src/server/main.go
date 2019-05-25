package main

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/service"
)

func main() {
	var tcpAddr string
	flag.StringVar(&tcpAddr, "tcp", "", "Run as a TCP server and listen on target address.")
	flag.Parse()
	if tcpAddr == "" {
		log.Fatalf("You must fill the TCP target address to listen.")
	}

	// Listen TCP target address
	log.Println("Starting tcp server.")
	lis, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("Failed to listen TCP: %v", err)
	}

	// For process
	stop := make(chan struct{})

	// Stop all listener by catching interrupt signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	go func(c chan os.Signal, lis net.Listener) {
		sig := <-c
		log.Printf("Caught signal: %s", sig.String())

		log.Printf("Stopping tcp listener...")
		lis.Close()

		log.Printf("TCP listener are stopped successfully.")
		close(stop)
	}(sigc, lis)

	// Init the service for handler to use
	sp := service.New()

	// Init the router
	router := NewRouter(sp)
	router.Handle(entity.GetBannersRequest_CMD, NewHandler(getBannersHandler))

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("Server|ListenerAccept|error:%v", err)
			break
		}

		go router.OnConnected(NewTCPConn(conn))
	}

	<-stop
	os.Exit(0)
}
