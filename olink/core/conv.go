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
		data, err := json.Marshal(msg)
		return data, err
	}
	return nil, nil
}

func (c *MessageConverter) FromData(data []byte) (Message, error) {
	switch c.Format {
	case FormatJson:
		var m Message
		err := json.Unmarshal(data, &m)
		return m, err
	}
	return nil, nil
}
