package core

import (
	"fmt"
)

type ProtocolListener interface {
	HandleLink(name string) error
	HandleUnlink(name string) error
	HandleInit(name string, props Props) error
	HandleSetProperty(name string, value Any) error
	HandlePropertyChange(name string, value Any) error
	HandleInvoke(id int, name string, args Args) error
	HandleInvokeReply(id int, name string, value Any) error
	HandleSignal(name string, args Args) error
	HandleError(msgType MsgType, id int, error string) error
}

type Protocol struct {
	listener ProtocolListener
	Base
}

func NewProtocol(listener ProtocolListener) *Protocol {
	return &Protocol{
		listener: listener,
		Base:     Base{},
	}
}

func NewLinkMessage(name string) Message {
	return Message{
		LINK,
		name,
	}
}

func NewInitMessage(name string, props Props) Message {
	return Message{
		INIT,
		name,
		props,
	}
}

func NewUnlinkMessage(name string) Message {
	return Message{
		UNLINK,
		name,
	}
}

func NewSetPropertyMessage(name string, value Any) Message {
	return Message{
		SET_PROPERTY,
		name,
		value,
	}
}

func NewPropertyChangeMessage(name string, value Any) Message {
	return Message{
		PROPERTY_CHANGE,
		name,
		value,
	}
}

func NewInvokeMessage(id int, name string, args Args) Message {
	return Message{
		INVOKE,
		id,
		name,
		args,
	}
}

func NewInvokeReplyMessage(id int, name string, value Any) Message {
	return Message{
		INVOKE_REPLY,
		id,
		name,
		value,
	}
}

func NewSignalMessage(name string, args Args) Message {
	return Message{
		SIGNAL,
		name,
		args,
	}
}

func NewErrorMessage(msgType MsgType, id int, error string) Message {
	return Message{
		msgType,
		id,
		error,
	}
}

func (p *Protocol) HandleMessage(msg Message) error {
	if p.listener == nil {
		return fmt.Errorf("no listener")
	}
	switch msg[0] {
	case LINK:
		return p.listener.HandleLink(msg[1].(string))
	case UNLINK:
		return p.listener.HandleUnlink(msg[1].(string))
	case INIT:
		return p.listener.HandleInit(msg[1].(string), msg[2].(Props))
	case SET_PROPERTY:
		return p.listener.HandleSetProperty(msg[1].(string), msg[2].(Any))
	case PROPERTY_CHANGE:
		return p.listener.HandlePropertyChange(msg[1].(string), msg[2].(Any))
	case INVOKE:
		return p.listener.HandleInvoke(msg[1].(int), msg[2].(string), msg[3].(Args))
	case INVOKE_REPLY:
		return p.listener.HandleInvokeReply(msg[1].(int), msg[2].(string), msg[3].(Any))
	case SIGNAL:
		return p.listener.HandleSignal(msg[1].(string), msg[2].(Args))
	case ERROR:
		return p.listener.HandleError(msg[0].(MsgType), msg[1].(int), msg[2].(string))
	}
	return fmt.Errorf("unknown message type: %v", msg)
}
