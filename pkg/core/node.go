package core

import (
	"errors"
	"fmt"
)

type BaseNode struct {
	writeFunc WriteMessageFunc
	converter MessageConverter
	protocol  Protocol
	Base
}

func NewBaseNode(writeFunc WriteMessageFunc) *BaseNode {
	return &BaseNode{writeFunc,
		MessageConverter{},
		Protocol{},
		Base{}}
}

func (node BaseNode) OnWrite(f WriteMessageFunc) {
	node.writeFunc = f
}

func (node BaseNode) EmitWrite(msg Message) {
	data, err := node.converter.ToString(msg)
	if err != nil {
		fmt.Printf("error converting message")
		return
	}
	if node.writeFunc != nil {
		node.writeFunc(data)
	}
}

func (node BaseNode) HandleMessage(data string) {
	msg, err := node.converter.FromString(data)
	if err != nil {
		fmt.Printf("error converting message")
		return
	}
	node.protocol.HandleMessage(msg)
}

func (node BaseNode) HandleLink(name string) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleUnlink(name string) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleInit(name string, props Props) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleSetProperty(name string, value Any) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandlePropertyChange(name string, value Any) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleInvoke(id int, name string, args Args) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleInvokeReply(id int, name string, value Any) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleSignal(name string, args Args) error {
	return errors.New("not implemented")
}

func (node BaseNode) HandleError(msgType MsgType, id int, error string) error {
	return errors.New("not implemented")
}
