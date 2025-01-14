package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

type Client struct {
	conn       net.Conn
	outbound   chan<- Command
	register   chan<- *Client
	unregister chan<- *Client
	username   string
}

func (c *Client) read() error {
	for {
		msg, err := bufio.NewReader(c.conn).ReadBytes('\n')
		if err == io.EOF {
			c.unregister <- c
			return nil
		}
		if err != nil {
			return err
		}

		c.handle(msg)

	}
}

func parseCommand(msg []byte) []byte {
	return bytes.ToUpper(bytes.TrimSpace(bytes.Split(msg, []byte(" "))[0]))
}

func parseArgs(msg []byte, cmd []byte) []byte {
	return bytes.TrimSpace(bytes.TrimPrefix(msg, cmd))
}

func (c *Client) handle(msg []byte) {
	cmd := parseCommand(msg)
	args := parseArgs(msg, cmd)

	switch string(cmd) {
	case "REG":
		if err := c.registerUser(args); err != nil {
			c.err(err)
		}
	case "JOIN":
		if err := c.joinChannel(args); err != nil {
			c.err(err)
		}
	case "LEAVE":
		if err := c.leaveChannel(args); err != nil {
			c.err(err)
		}
	case "MSG":
		if err := c.sendMsg(args); err != nil {
			c.err(err)
		}
	case "CHNS":
		c.getChannels()
	case "USRS":
		c.getUsers()
	default:
		c.err(fmt.Errorf("unknown command: %s", cmd))
	}
}

func (c *Client) registerUser(args []byte) error {
	user := bytes.TrimSpace(args)
	if user[0] != '@' {
		return fmt.Errorf("Username must start with '@'")
	}
	if len(user) == 0 {
		return fmt.Errorf("Username cannot be empty")
	}

	c.username = string(user)
	c.register <- c

	return nil

}

func (c *Client) err(err error) {
	c.conn.Write([]byte(fmt.Sprintf("ERR: %s\n", err.Error())))
}
