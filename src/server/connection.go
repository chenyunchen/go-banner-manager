package main

import (
	"encoding/binary"
	"io"
	"time"
)

// Packet represents a general data packet
type Packet struct {
	Content []byte
	Conn    Conn
}

// Conn is the interface of a general connection can write message
type Conn interface {
	LocalAddr() string
	RemoteAddr() string
	io.Writer
	WriteMsg(msg string, byteOrder binary.ByteOrder, timeout time.Duration) error
	Close() error
}

// EventHandler is the interface to handle event
type EventHandler interface {
	OnConnected(conn Conn)
	OnPacket(packet *Packet)
	OnDisconnected(conn Conn)
}
