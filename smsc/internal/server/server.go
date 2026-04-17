package server

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Server struct {
	addr         string
	onConnect    func(conn net.Conn)
	onDisconnect func(conn net.Conn)
	handler      func(conn net.Conn, raw []byte)
}

func NewServer(addr string, onConnect func(net.Conn), onDisconnect func(net.Conn), handler func(net.Conn, []byte)) *Server {
	return &Server{
		addr:         addr,
		onConnect:    onConnect,
		onDisconnect: onDisconnect,
		handler:      handler,
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()
	log.Println("Listening on", s.addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	s.onConnect(conn)
	defer func() {
		s.onDisconnect(conn)
		conn.Close()
	}()

	for {
		lenBuf := make([]byte, 4)
		if _, err := io.ReadFull(conn, lenBuf); err != nil {
			log.Println("Connection closed:", err)
			return
		}

		commandLength := binary.BigEndian.Uint32(lenBuf)
		if commandLength < 16 {
			log.Println("Malformed PDU: command_length too small:", commandLength)
			return
		}

		commandBuf := make([]byte, commandLength-4)
		if _, err := io.ReadFull(conn, commandBuf); err != nil {
			log.Println("Error reading PDU body:", err)
			return
		}

		raw := append(lenBuf, commandBuf...)
		s.handler(conn, raw)
	}
}
