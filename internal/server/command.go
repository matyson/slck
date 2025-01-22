package server

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
