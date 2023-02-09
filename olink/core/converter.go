package core

import (
	"bytes"
	"encoding/json"
)

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

func NewConverter(format MessageFormat) *MessageConverter {
	return &MessageConverter{
		Format: format,
	}
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
		var msg Message
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		err := decoder.Decode(&msg)
		return msg, err
	}
	return nil, nil
}
