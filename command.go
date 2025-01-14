package main

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
	ID        ID
	body      []byte
	recipient string
	sender    string
}
