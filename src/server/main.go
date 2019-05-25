package main

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
	"bytes"
	"encoding/binary"
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

		go handleConnection(conn, router)
	}

	<-stop
	os.Exit(0)
}

// TODO: move tcp to outside
func handleConnection(conn net.Conn, router *router) {
	for {
		var (
			contentSize uint32
			isHead      bool = true
			buffer           = bytes.NewBuffer(make([]byte, 0, 200))
			bytes            = make([]byte, 200)
			head             = make([]byte, 4)
			content          = make([]byte, 200)
		)
		readLen, err := conn.Read(bytes)
		if err != nil {
			return
		}
		_, err = buffer.Write(bytes[0:readLen])
		if err != nil {
			log.Printf("Server|BufferWrite|error:%v", err)
			return
		}
		for {
			if isHead {
				if buffer.Len() >= 2 {
					_, err := buffer.Read(head)
					if err != nil {
						log.Printf("Server|BufferRead|error:%v", err)
						return
					}
					contentSize = binary.BigEndian.Uint32(head)
					isHead = false
				} else {
					break
				}
			}
			if !isHead {
				if buffer.Len() >= int(contentSize) {
					_, err := buffer.Read(content[:contentSize])
					if err != nil {
						log.Printf("Server|BufferRead|error:%v", err)
						return
					}
					input := NewTCPPacket(content[:contentSize], conn)
					router.OnPacket(&input)
					isHead = true
				} else {
					break
				}
			}

		}
	}
}
