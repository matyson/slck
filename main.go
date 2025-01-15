package main

import "net"

func main() {
	hub := &Hub{
		channels:        make(map[string]*Channel),
		clients:         make(map[string]*Client),
		commands:        make(chan Command),
		unregistrations: make(chan *Client),
		registrations:   make(chan *Client)}

	go hub.run()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		client := &Client{
			conn:       conn,
			outbound:   hub.commands,
			register:   hub.registrations,
			unregister: hub.unregistrations}

		go client.read()
	}

}
