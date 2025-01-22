package main

import (
	"log"
	"net"

	"github.com/matyson/slck/internal/server"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("%v", err)
	}

	hub := server.NewHub()
	go hub.Run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%v", err)
		}

		c := server.NewClient(conn, hub.Commands, hub.Registrations, hub.Unresgistrations)
		go c.Read()
	}
}
