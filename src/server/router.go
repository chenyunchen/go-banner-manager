package main

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/service"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

const (
	BUFFER_SIZE = 1024
	HEAD_SIZE   = 4
)

type fn func(*context)

// context contains service and input from client side
type context struct {
	service *service.Service
	input   *Packet
}

// NewContext creates a new router
func NewContext(sp *service.Service, input *Packet) *context {
	return &context{
		service: sp,
		input:   input,
	}
}

// router contains cmd mapping function and handle connection
type router struct {
	handlers    map[uint16]fn
	service     *service.Service
	onConnected func(conn Conn)
}

// NewRouter creates a new router
func NewRouter(sp *service.Service) *router {
	return &router{
		service:  sp,
		handlers: make(map[uint16]fn),
	}
}

// NewRouter creates a new handler
func NewHandler(f func(*context)) fn {
	return fn(f)
}

// Handle register a handler for a cmd
func (r *router) Handle(cmd uint16, fn fn) error {
	_, present := r.handlers[cmd]
	if present {
		return fmt.Errorf("Server|Router|Handle|CMDOccupied|cmd:%d", cmd)
	}
	r.handlers[cmd] = fn
	return nil
}

// OnConnected implements the EventHandler interface
func (r *router) OnConnected(conn Conn) {
	var (
		contentSize uint32
		isHead      bool = true
		buffer           = bytes.NewBuffer(make([]byte, 0, BUFFER_SIZE))
		bytes            = make([]byte, BUFFER_SIZE)
		head             = make([]byte, HEAD_SIZE)
		content          = make([]byte, BUFFER_SIZE)
	)
	for {
		readLen, err := conn.Read(bytes)
		if err != nil {
			log.Printf("Server|Router|ConnRead|error:%v", err)
			return
		}
		_, err = buffer.Write(bytes[0:readLen])
		if err != nil {
			log.Printf("Server|Router|BufferWrite|error:%v", err)
			return
		}
		for {
			if isHead {
				if buffer.Len() >= HEAD_SIZE {
					_, err := buffer.Read(head)
					if err != nil {
						log.Printf("Server|Router|BufferRead|error:%v", err)
						return
					}
					contentSize = binary.BigEndian.Uint32(head)
					isHead = false
				} else {
					break
				}
			}
			if !isHead {
				log.Printf("Packet Receive: %v bytes\n", buffer.Len())
				log.Printf("TCP Shoud Send: %v bytes\n", contentSize)
				if buffer.Len() >= int(contentSize) {
					_, err := buffer.Read(content[:contentSize])
					if err != nil {
						log.Printf("Server|Router|BufferRead|error:%v", err)
						return
					}
					input := Packet{
						Content: content[:contentSize],
						Conn:    conn,
					}
					r.OnPacket(&input)
					isHead = true
				} else {
					break
				}
			}

		}
	}
}

// OnPacket implements the EventHandler interface
func (r *router) OnPacket(packet *Packet) {
	if len(packet.Content) < 2 {
		log.Printf("Server|Router|OnPacket|InvalidPacket")
		return
	}

	cmd := binary.BigEndian.Uint16([]byte(packet.Content))
	packet.Content = packet.Content[2:]

	r.handlers[cmd](NewContext(r.service, packet))

	return
}
