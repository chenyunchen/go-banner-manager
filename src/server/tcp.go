package main

import (
	"encoding/binary"
	"net"
	"time"
)

// tcpConn contains a net.Conn
type tcpConn struct {
	conn net.Conn
}

// LocalAddr implements the Conn interface
func (c tcpConn) LocalAddr() string {
	return c.conn.LocalAddr().(*net.TCPAddr).IP.String() // 127.0.0.1
}

// RemoteAddr implements the Conn interface
func (c tcpConn) RemoteAddr() string {
	return c.conn.RemoteAddr().(*net.TCPAddr).IP.String() // 127.0.0.1
}

// Read implements the io.Read
func (c tcpConn) Read(buf []byte) (n int, err error) {
	return c.conn.Read(buf)
}

// WriteMsg implements the Conn interface
func (c tcpConn) WriteMsg(msg string, byteOrder binary.ByteOrder, timeout time.Duration) error {
	_, err := c.conn.Write([]byte(msg))
	return err
}

// Write implements the io.Writer
func (c tcpConn) Write(buf []byte) (n int, err error) {
	return c.conn.Write(buf)
}

// Close implements the Conn interface
func (c tcpConn) Close() error {
	return c.conn.Close()
}

// NewTCPConn will new a tcpConn from net.Conn
func NewTCPConn(conn net.Conn) tcpConn {
	return tcpConn{
		conn: conn,
	}
}

func NewTCPPacket(content []byte, conn net.Conn) Packet {
	return Packet{
		Content: content,
		Conn:    NewTCPConn(conn),
	}
}
