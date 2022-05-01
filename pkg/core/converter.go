package core

import "encoding/json"

type MessageFormat int

const (
	FormatJson    MessageFormat = 1
	FormatBson    MessageFormat = 2
	FormatMsgPack MessageFormat = 3
	FormatCbor    MessageFormat = 4
)

type MessageConverter struct {
	Format MessageFormat
}

func (c *MessageConverter) ToData(msg Message) ([]byte, error) {
	switch c.Format {
	case FormatJson:
		b, err := json.Marshal(msg)
		return b, err
	}
	return nil, nil
}

func (c *MessageConverter) FromData(msg []byte) (Message, error) {
	switch c.Format {
	case FormatJson:
		var m Message
		err := json.Unmarshal(msg, &m)
		return m, err
	}
	return nil, nil
}
