package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MsgType int

const (
	LINK            MsgType = 10
	INIT            MsgType = 11
	UNLINK          MsgType = 12
	SET_PROPERTY    MsgType = 20
	PROPERTY_CHANGE MsgType = 21
	INVOKE          MsgType = 30
	INVOKE_REPLY    MsgType = 31
	SIGNAL          MsgType = 40
	ERROR           MsgType = 90
)

type MessageFormat int

const (
	JSON    MessageFormat = 1
	BSON    MessageFormat = 2
	MSGPACK MessageFormat = 3
	CBOR    MessageFormat = 4
)

// name: <resource>/<path>
// resource: <module>.<interface>
func ResourceFromName(name string) string {
	return strings.Split(name, "/")[0]
}

func PathFromName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

func HasPath(name string) bool {
	return strings.Contains(name, "/")
}

func CreateName(resource, path string) string {
	return fmt.Sprintf("%s/%s", resource, path)
}

type WriteMessageFunc func(msg string)

type MessageWriter interface {
	WriteMessage(msg interface{}) error
}

type MessageHandlers interface {
	HandleMessage(data string)
}

type Base struct {
}

type Message []interface{}
type Args []interface{}
type Props map[string]interface{}
type Any interface{}

type MessageConverter struct {
	format MessageFormat
}

func (c *MessageConverter) ToString(msg Message) (string, error) {
	switch c.format {
	case JSON:
		b, err := json.Marshal(msg)
		return string(b), err
	}
	return "", nil
}

func (c *MessageConverter) FromString(msg string) (Message, error) {
	switch c.format {
	case JSON:
		var m Message
		err := json.Unmarshal([]byte(msg), &m)
		return m, err
	}
	return nil, nil
}
