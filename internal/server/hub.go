package server

import (
	"fmt"
	"strings"
)

type Hub struct {
	channels        map[string]*Channel
	clients         map[string]*Client
	commands        chan Command
	unregistrations chan *Client
	registrations   chan *Client
}

func newHub() *Hub {
	return &Hub{
		channels:        make(map[string]*Channel),
		clients:         make(map[string]*Client),
		commands:        make(chan Command),
		unregistrations: make(chan *Client),
		registrations:   make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.registrations:
			h.register(client)
		case client := <-h.unregistrations:
			h.unregister(client)
		case command := <-h.commands:
			switch command.id {
			case JOIN:
				h.joinChannel(command.sender, command.recipient)
			case LEAVE:
				h.leaveChannel(command.sender, command.recipient)
			case MSG:
				h.sendMessage(command.sender, command.recipient, command.body)
			case CHNS:
				h.listChannels(command.sender)
			case USRS:
				h.listUsers(command.sender)
			default:
				h.err(fmt.Errorf("unknown command: %d", command.id))
			}

		}
	}
}

func (h *Hub) register(client *Client) {
	_, ok := h.clients[client.username]
	if ok {
		client.conn.Write([]byte("ERR: Username already taken\n"))
	} else {
		h.clients[client.username] = client
		client.conn.Write([]byte(fmt.Sprintf("OK: Welcome %s\n", client.username)))
	}
}

func (h *Hub) unregister(client *Client) {
	_, ok := h.clients[client.username]
	if ok {
		delete(h.clients, client.username)
		for _, channel := range h.channels {
			delete(channel.clients, client)
		}
		client.conn.Write([]byte(fmt.Sprintf("OK: Goodbye %s\n", client.username)))
	} else {
		client.conn.Write([]byte("ERR: Username not found\n"))
	}
}

func (h *Hub) joinChannel(user string, channelName string) {
	client, ok := h.clients[user]
	if ok {
		channel, ok := h.channels[channelName]
		if ok {
			channel.clients[client] = true
		} else {
			h.channels[channelName] = &Channel{
				name:    channelName,
				clients: make(map[*Client]bool),
			}
			h.channels[channelName].clients[client] = true
		}
	}
}

func (h *Hub) leaveChannel(user string, channelName string) {
	client, ok := h.clients[user]
	if ok {
		channel, ok := h.channels[channelName]
		if ok {
			delete(channel.clients, client)
		}
	}
}

func (h *Hub) sendMessage(user string, recipient string, message []byte) {
	if sender, ok := h.clients[user]; ok {
		switch recipient[0] {
		case '@':
			if user, ok := h.clients[recipient]; ok {
				msg := append([]byte(fmt.Sprintf("%s: ", sender.username)), message...)
				msg = append(msg, '\n')

				user.conn.Write(msg)
			} else {
				sender.conn.Write([]byte("ERR: User not found\n"))
			}
		case '#':
			if channel, ok := h.channels[recipient]; ok {
				if _, ok := channel.clients[sender]; ok {
					channel.broadcast(sender.username, message)
				}
			} else {
				sender.conn.Write([]byte("ERR: Channel not found\n"))
			}
		default:
			sender.conn.Write([]byte("ERR: Recipient must be a #channel or @user\n"))
		}
	}

}

func (h *Hub) listChannels(user string) {
	if client, ok := h.clients[user]; ok {
		var names []string

		if len(h.channels) == 0 {
			client.conn.Write([]byte("No channels\n"))
		} else {
			for name := range h.channels {
				names = append(names, name)
			}

			res := strings.Join(names, ", ")
			client.conn.Write([]byte(fmt.Sprintf("Channels: %s\n", res)))
		}
	}
}

func (h *Hub) listUsers(user string) {
	if client, ok := h.clients[user]; ok {
		var names []string

		if len(h.clients) == 0 {
			client.conn.Write([]byte("No users\n"))
		} else {
			for name := range h.clients {
				names = append(names, name)
			}

			res := strings.Join(names, ", ")
			client.conn.Write([]byte(fmt.Sprintf("Users: %s\n", res)))
		}
	}
}

func (h *Hub) err(err error) {
	fmt.Printf("Error: %s\n", err.Error())
}
