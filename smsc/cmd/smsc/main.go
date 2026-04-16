package main

import (
	"log"
	"net"

	"github.com/joeygalvin/smpp-sandbox/smsc/internal/server"
)

func main() {
	s := server.NewServer(":2775", func(conn net.Conn, raw []byte) {
		log.Println("Received data:", string(raw))
	})
	s.Start()
}
