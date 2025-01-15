package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
)

var DELIMITER = []byte("\r\n")

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

func (c *Client) joinChannel(args []byte) error {
	channel := bytes.TrimSpace(args)
	if channel[0] != '#' {
		return fmt.Errorf("Channel must start with '#'")
	}

	c.outbound <- Command{
		recipient: string(channel),
		sender:    c.username,
		id:        JOIN,
	}

	return nil
}

func (c *Client) leaveChannel(args []byte) error {
	channel := bytes.TrimSpace(args)
	if channel[0] != '#' {
		return fmt.Errorf("Channel must start with '#'")
	}

	c.outbound <- Command{
		recipient: string(channel),
		sender:    c.username,
		id:        LEAVE,
	}

	return nil
}

func (c *Client) sendMsg(args []byte) error {
	// MSG #channel or @user lenght\r\message
	args = bytes.TrimSpace(args)
	if args[0] != '#' && args[0] != '@' {
		return fmt.Errorf("Recipient must be a #channel or @user")
	}

	recipient := bytes.Split(args, []byte(" "))[0]
	if len(recipient) == 0 {
		return fmt.Errorf("Recipient must have a name")
	}

	args = bytes.TrimSpace(bytes.TrimPrefix(args, recipient))
	l := bytes.Split(args, DELIMITER)[0]
	length, err := strconv.Atoi(string(l))
	if err != nil {
		return fmt.Errorf("body length must be present")

	}
	if length == 0 {
		return fmt.Errorf("body length must be at least 1")
	}

	padding := len(l) + len(DELIMITER) // Size of the body length + the delimiter
	body := args[padding : padding+length]

	c.outbound <- Command{
		recipient: string(recipient),
		sender:    c.username,
		body:      body,
		id:        MSG,
	}

	return nil
}

func (c *Client) getChannels() {
	c.outbound <- Command{
		sender: c.username,
		id:     CHNS,
	}
}

func (c *Client) getUsers() {
	c.outbound <- Command{
		sender: c.username,
		id:     USRS,
	}
}
