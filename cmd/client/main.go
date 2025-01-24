package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

type ID int

const (
	REG ID = iota
	JOIN
	LEAVE
	MSG
	CHNS
	USRS
)

type Command struct {
	id        ID
	body      []byte
	recipient string
	sender    string
}

func main() {

	cmd := Command{id: REG, body: []byte("@matyson")}

	r := newRequest("localhost:3000", cmd)
	err := r.send()
	if err != nil {
		fmt.Println(err)
	}

}

type Request struct {
	addr    string
	command Command
}

func newRequest(addr string, command Command) *Request {
	return &Request{
		addr:    addr,
		command: command,
	}
}

func (r *Request) send() error {
	conn, err := net.Dial("tcp", r.addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	switch r.command.id {
	case REG:
		r.register(conn)
	case JOIN:
		// send join command
	case LEAVE:
		// send leave command
	case MSG:
		// send message command
	case CHNS:
		// send channels command
	case USRS:
	// send users command
	default:
		return fmt.Errorf("unknown command ID: %d", r.command.id)
	}
	return nil
}

func (r *Request) register(conn net.Conn) error {
	user := r.command.body
	if user[0] != '@' {
		return fmt.Errorf("username must start with '@'")
	}
	fmt.Println("sending registration command for ", string(user))
	_, err := conn.Write([]byte("REG " + string(user) + "\n"))
	if err != nil {
		return err
	}
	for {
		msg, err := bufio.NewReader(conn).ReadBytes('\n')
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		fmt.Print(string(msg))
		status := bytes.Split(msg, []byte(" "))[0]

		if string(status) == "OK:" {
			fmt.Println("registration successful")
			return nil
		}
		if string(status) == "ERR:" {
			fmt.Println("registration failed")
			return nil
		}

	}
}
