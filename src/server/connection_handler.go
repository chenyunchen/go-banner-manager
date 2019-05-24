package main

import (
	"net"
)

func (s *server) handleConnection(conn net.Conn) {
	defer conn.Close()
}
