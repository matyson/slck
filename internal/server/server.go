package server

import (
	"fmt"
	"log"
	"net"
)

type server struct {
	Ln   net.Listener
	hub  *Hub
	addr string
}

func (s *server) AcceptLoop() error {
	for {
		conn, err := s.Ln.Accept()
		if err != nil {
			return err
		}
		log.Printf("new connection from %s\n", conn.RemoteAddr())
		conn.Write([]byte("Welcome to SLCK\n"))

		c := newClient(conn, s.hub.commands, s.hub.registrations, s.hub.unregistrations)
		go c.read()

	}
}

func NewServer(addr string) *server {
	return &server{
		hub:  newHub(),
		addr: addr,
	}
}

func (s *server) Run() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Println(err)
	}
	s.Ln = ln
	log.Printf("server listening on %s\n", s.addr)

	go s.hub.run()

}
