package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("%v", err)
	}

	hub := &Hub{
		channels:        make(map[string]*Channel),
		clients:         make(map[string]*Client),
		commands:        make(chan Command),
		unregistrations: make(chan *Client),
		registrations:   make(chan *Client)}
	go hub.run()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%v", err)
		}

		c := &Client{
			conn:       conn,
			outbound:   hub.commands,
			register:   hub.registrations,
			unregister: hub.unregistrations}

		go c.read()
	}
}
