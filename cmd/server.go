package main

import (
	"log"

	"github.com/matyson/slck/internal/server"
)

const (
	PORT = "3000"
)

func main() {
	s := server.NewServer(":" + PORT)
	s.Run()
	log.Fatal(s.AcceptLoop())

	defer s.Ln.Close()

}
