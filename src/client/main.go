package main

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"time"
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

	// TODO: Write Test
	data := entity.UpdateBannerRequest{
		Serial:      1,
		StartedTime: uint32(time.Now().Unix()),
		ExpiredTime: uint32(time.Now().Unix()),
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.GetBannersRequest_CMD))
	conn.Write(append(output, b...))

	// TODO: Read Test
	scannerConn := bufio.NewScanner(conn)
	for scannerConn.Scan() {
		log.Println("Server sends: " + scannerConn.Text())
		break
	}
}
