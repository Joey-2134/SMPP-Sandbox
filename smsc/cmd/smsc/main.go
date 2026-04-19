package main

import (
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeygalvin/smpp-sandbox/smsc/internal/server"
	"github.com/joeygalvin/smpp-sandbox/smsc/internal/session"
	"github.com/joeygalvin/smpp-sandbox/smsc/internal/store"
	"github.com/joeygalvin/smpp-sandbox/smsc/tui"
)

func main() {
	store, err := store.Open(store.DefaultPath)
	if err != nil {
		log.Fatalf("Failed to open store: %v", err)
	}
	defer store.Close()
	manager := session.NewManager(store)

	s := server.NewServer(
		":2775",
		func(conn net.Conn) {
			manager.AddSession(conn)
			log.Println("Client connected:", conn.RemoteAddr())
		},
		func(conn net.Conn) {
			manager.RemoveSession(conn)
			log.Println("Client disconnected:", conn.RemoteAddr())
		},
		func(conn net.Conn, raw []byte) {
			sess := manager.GetSession(conn)
			if err := sess.Handle(raw); err != nil {
				log.Println("Session error:", err)
			}
		},
	)

	go s.Start()

	p := tea.NewProgram(tui.NewModel(store), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
