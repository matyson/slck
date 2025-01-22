package main

import (
	"log"
	"net"

	"github.com/matyson/slck/internal/server"
)

const (
	PORT = "8080"
)

func main() {
	ln, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Printf("%v", err)
	}
	defer ln.Close()
	log.Printf("server listening on :%s\n", PORT)

	hub := server.NewHub()
	go hub.Run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%v", err)
		}
		log.Printf("new connection from %s\n", conn.RemoteAddr())

		c := server.NewClient(conn, hub.Commands, hub.Registrations, hub.Unresgistrations)
		go c.Read()
	}
}
