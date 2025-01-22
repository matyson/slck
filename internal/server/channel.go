package server

type Channel struct {
	name    string // start with #
	clients map[*Client]bool
}

func (c *Channel) broadcast(subject string, m []byte) {
	msg := append([]byte(subject), ": "...)
	msg = append(msg, m...)
	msg = append(msg, '\n')

	for client := range c.clients {
		client.conn.Write(msg)
	}
}
