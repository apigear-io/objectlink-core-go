package core

type MessageWriter interface {
	WriteMessage(m Message) error
}
