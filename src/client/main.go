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
	"strconv"
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
			log.Fatalf("You should fill the serial number, start time and expire time.")
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
			log.Fatalf("Please fill the corret start time number.")
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
		expiredTime, err := strconv.ParseUint(start, 10, 64)
		if err != nil {
			log.Fatalf("Please fill the corret expire time number.")
		}

		updateBannerExpiredTime(conn, uint16(serial), uint32(expiredTime))
	case "clear_all_timers":
		clearAllBannerTimers(conn)
	}

	scannerConn := bufio.NewScanner(conn)
	for scannerConn.Scan() {
		banners := []entity.Banner{}
		json.Unmarshal(scannerConn.Bytes(), &banners)
		log.Println("Display Banner:")
		for _, banner := range banners {
			log.Printf("Serial: %d\n", banner.Serial)
			log.Printf("Event: %s\n", banner.Event)
			log.Printf("Text: %s\n", banner.Text)
			log.Printf("Image: %s\n", banner.Image)
			log.Printf("URL: %s\n", banner.URL)
			log.Printf("Started Time: %s\n", banner.StartedTime)
			log.Printf("Expired Time: %s\n", banner.ExpiredTime)
		}
		break
	}
}

func getBanners(conn net.Conn) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2))
	conn.Write(size)

	cmd := make([]byte, 2)
	binary.BigEndian.PutUint16(cmd, uint16(entity.GetBannersRequest_CMD))
	conn.Write(cmd)
}

func updateBanner(conn net.Conn, serial uint16, startedTime, expiredTime uint32) {
	data := entity.UpdateBannerRequest{
		Serial:      serial,
		StartedTime: startedTime,
		ExpiredTime: expiredTime,
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.UpdateBannerRequest_CMD))
	conn.Write(append(output, b...))
}

func updateBannerStartedTime(conn net.Conn, serial uint16, startedTime uint32) {
	data := entity.UpdateBannerStartedTimeRequest{
		Serial:      serial,
		StartedTime: startedTime,
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.UpdateBannerStartedTimeRequest_CMD))
	conn.Write(append(output, b...))
}

func updateBannerExpiredTime(conn net.Conn, serial uint16, expiredTime uint32) {
	data := entity.UpdateBannerExpiredTimeRequest{
		Serial:      serial,
		ExpiredTime: expiredTime,
	}
	b, _ := json.Marshal(data)

	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2+len(b)))
	conn.Write(size)

	output := make([]byte, 2)
	binary.BigEndian.PutUint16(output, uint16(entity.UpdateBannerExpiredTimeRequest_CMD))
	conn.Write(append(output, b...))
}

func clearAllBannerTimers(conn net.Conn) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(2))
	conn.Write(size)

	cmd := make([]byte, 2)
	binary.BigEndian.PutUint16(cmd, uint16(entity.ClearAllBannerTimersRequest_CMD))
	conn.Write(cmd)
}
