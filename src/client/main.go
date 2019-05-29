package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/entity"
)

func main() {
	var action string
	var tcpAddr string
	var serial string
	var start string
	var expire string
	flag.StringVar(&tcpAddr, "tcp", "0.0.0.0:8080", "Run as a TCP client to connect target address.")
	flag.StringVar(&action, "action", "get", "Run a command to control the banner-manager.")
	flag.StringVar(&serial, "serial", "", "Decide which the banner display. (Serial Number)")
	flag.StringVar(&start, "start", "", "Decide when the banner display. (Unix Timestamp)")
	flag.StringVar(&expire, "expire", "", "Decide when the banner expire. (Unix Timestamp)")
	flag.Parse()

	// Connect TCP target address
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		log.Printf("Dial failed: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	switch action {
	case "get":
		getBanners(conn)
	case "update":
		if serial == "" || start == "" || expire == "" {

		}
		serial, err := strconv.ParseUint(serial, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret serial number.")
		}
		startedTime, err := strconv.ParseUint(start, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret start time number.")
		}
		expiredTime, err := strconv.ParseUint(expire, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret expire time number.")
		}

		updateBanner(conn, uint16(serial), uint32(startedTime), uint32(expiredTime))
	case "update_start":
		if serial == "" || start == "" {
			log.Fatalf("You should fill the serial number and start time.")
		}
		serial, err := strconv.ParseUint(serial, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret serial number.")
		}
		startedTime, err := strconv.ParseUint(start, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret start time number.")
		}

		updateBannerStartedTime(conn, uint16(serial), uint32(startedTime))
	case "update_expire":
		if serial == "" || expire == "" {
			log.Fatalf("You should fill the serial number and expire time.")
		}
		serial, err := strconv.ParseUint(serial, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret serial number.")
		}
		expiredTime, err := strconv.ParseUint(expire, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret expire time number.")
		}

		updateBannerExpiredTime(conn, uint16(serial), uint32(expiredTime))
	case "clear_all_timers":
		clearAllBannerTimers(conn)
	}

	scannerConn := bufio.NewScanner(conn)
	for scannerConn.Scan() {
		b := scannerConn.Bytes()
		// TODO: Not using bytes length to judge if the error occur. This will fail if the origin data's size is small.
		if len(b) <= 100 {
			log.Println(string(b))
			break
		}
		banner := entity.Banner{}
		json.Unmarshal(b, &banner)
		log.Println("Display Banner:")
		log.Printf("Serial: %d\n", banner.Serial)
		log.Printf("Event: %s\n", banner.Event)
		log.Printf("Text: %s\n", banner.Text)
		log.Printf("Image: %s\n", banner.Image)
		log.Printf("URL: %s\n", banner.URL)
		log.Printf("Started Time: %s\n", banner.StartedTime)
		log.Printf("Expired Time: %s\n", banner.ExpiredTime)
		break
	}
}
