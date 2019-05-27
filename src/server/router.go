package main

import (
	"4670e1812919d92b8cf4e33ac38bc40e449521da/src/service"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type fn func(*context)

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

type router struct {
	handlers       map[uint16]fn
	service        *service.Service
	onConnected    func(conn Conn)
	onDisConnected func(conn Conn)
}

// NewRouter creates a new router
func NewRouter(sp *service.Service) *router {
	return &router{
		service:  sp,
		handlers: make(map[uint16]fn),
	}
}

func NewHandler(f func(*context)) fn {
	return fn(f)
}

// Handle register a handler for a cmd
func (r *router) Handle(cmd uint16, fn fn) error {
	_, present := r.handlers[cmd]
	if present {
		return fmt.Errorf("Router|Handle|CMDOccupied|cmd:%d", cmd)
	}
	r.handlers[cmd] = fn
	return nil
}

func (r *router) HandleConnected(f func(conn Conn)) error {
	if r.onConnected != nil {
		return errors.New("Router|HandleConnected|ConnectHandlerExists")
	}
	r.onConnected = f
	return nil
}

// OnConnected implements the EventHandler interface
func (r *router) OnConnected(conn Conn) {
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
				if buffer.Len() >= 4 {
					_, err := buffer.Read(head)
					if err != nil {
						log.Printf("Server|BufferRead|error:%v", err)
						return
					}
					contentSize = binary.BigEndian.Uint32(head)
					isHead = false
				} else {
					return
				}
			}
			if !isHead {
				if buffer.Len() >= int(contentSize) {
					_, err := buffer.Read(content[:contentSize])
					if err != nil {
						log.Printf("Server|BufferRead|error:%v", err)
						return
					}
					input := Packet{
						Content: content[:contentSize],
						Conn:    conn,
					}
					r.OnPacket(&input)
					isHead = true
				} else {
					return
				}
			}

		}
	}
}

// OnDisconnected implements the EventHandler interface
func (r *router) OnDisconnected(conn Conn) {
	if r.onDisConnected != nil {
		r.onDisConnected(conn)
	}

	return
}

// OnPacket implements the EventHandler interface
func (r *router) OnPacket(packet *Packet) {
	if len(packet.Content) < 2 {
		log.Printf("Router|OnPacket|InvalidPacket")
		return
	}

	cmd := binary.BigEndian.Uint16([]byte(packet.Content))
	packet.Content = packet.Content[2:]

	r.handlers[cmd](NewContext(r.service, packet))

	return
}
