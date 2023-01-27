package core

import (
	"fmt"
)

type MsgType int64

func (t MsgType) String() string {
	switch t {
	case MsgUnknown:
		return "unknown"
	case MsgLink:
		return "link"
	case MsgInit:
		return "init"
	case MsgUnlink:
		return "unlink"
	case MsgSetProperty:
		return "set"
	case MsgPropertyChange:
		return "change"
	case MsgInvoke:
		return "invoke"
	case MsgInvokeReply:
		return "reply"
	case MsgSignal:
		return "signal"
	case MsgError:
		return "error"
	}
	return fmt.Sprintf("unknown(%d)", t)
}

func MsgTypeFromString(s string) MsgType {
	switch s {
	case "unknown":
		return MsgUnknown
	case "link":
		return MsgLink
	case "init":
		return MsgInit
	case "unlink":
		return MsgUnlink
	case "set":
		return MsgSetProperty
	case "change":
		return MsgPropertyChange
	case "invoke":
		return MsgInvoke
	case "reply":
		return MsgInvokeReply
	case "signal":
		return MsgSignal
	case "error":
		return MsgError
	}
	return MsgUnknown
}

const (
	MsgUnknown        MsgType = 0
	MsgLink           MsgType = 10
	MsgInit           MsgType = 11
	MsgUnlink         MsgType = 12
	MsgSetProperty    MsgType = 20
	MsgPropertyChange MsgType = 21
	MsgInvoke         MsgType = 30
	MsgInvokeReply    MsgType = 31
	MsgSignal         MsgType = 40
	MsgError          MsgType = 90
)

type Args []any

func AsKwArgs(a Args) KWArgs {
	kwargs := make(KWArgs)
	for i, v := range a {
		name := fmt.Sprintf("arg%d", i)
		kwargs[name] = v
	}
	return kwargs
}

type KWArgs map[string]any

func (a KWArgs) Keys() []string {
	keys := make([]string, len(a))
	i := 0
	for k := range a {
		keys[i] = k
		i++
	}
	return keys
}

type Any any

type Message []any

func (m Message) Type() MsgType {
	return AsMsgType(m[0])
}

func (m Message) StringType() string {
	return AsString(m[0])
}

// AsInit returns the name and props of the init message
func (m Message) AsInit() (string, KWArgs) {
	return AsString(m[1]), AsProps(m[2])
}

// AsLink returns the name of the link message
func (m Message) AsLink() string {
	return AsString(m[1])
}

// AsUnlink returns the name of the link message
func (m Message) AsUnlink() string {
	return AsString(m[1])
}

// AsSetProperty returns the name and value of the message
// message := MsgType, PropertyId, Value
func (m Message) AsSetProperty() (string, Any) {
	return AsString(m[1]), AsAny(m[2])
}

// AsPropertyChange returns the name and value of the property change
// message := MsgType, PropertyId, Value
func (m Message) AsPropertyChange() (string, Any) {
	return AsString(m[1]), AsAny(m[2])
}

// AsInvoke returns the id, name and args of the invoke message
// message := MsgType, RequestId, MethodId, Args
func (m Message) AsInvoke() (int64, string, Args) {
	return AsInt(m[1]), AsString(m[2]), AsArgs(m[3])
}

// AsInvokeReply returns the id and result of the invoke reply
// message := MsgType, RequestId, MethodId, Value
func (m Message) AsInvokeReply() (int64, string, Any) {
	return AsInt(m[1]), AsString(m[2]), AsAny(m[3])
}

// AsSignal returns the name and args of the signal message
// message := MsgType, SignalId, Args
func (m Message) AsSignal() (string, Args) {
	return AsString(m[1]), AsArgs(m[2])
}

func (m Message) AsError() (MsgType, int64, string) {
	return AsMsgType(m[0]), AsInt(m[1]), AsString(m[2])
}

func MakeLinkMessage(objectId string) Message {
	return Message{
		MsgLink,
		objectId,
	}
}

func MakeInitMessage(objectId string, props KWArgs) Message {
	return Message{
		MsgInit,
		objectId,
		props,
	}
}

func MakeUnlinkMessage(objectId string) Message {
	return Message{
		MsgUnlink,
		objectId,
	}
}

func MakeSetPropertyMessage(propertyId string, value Any) Message {
	return Message{
		MsgSetProperty,
		propertyId,
		value,
	}
}

func MakePropertyChangeMessage(propertyId string, value Any) Message {
	return Message{
		MsgPropertyChange,
		propertyId,
		value,
	}
}

func MakeInvokeMessage(requestId int64, methodId string, args Args) Message {
	return Message{
		MsgInvoke,
		requestId,
		methodId,
		args,
	}
}

func MakeInvokeReplyMessage(requestId int64, methodId string, value Any) Message {
	return Message{
		MsgInvokeReply,
		requestId,
		methodId,
		value,
	}
}

func MakeSignalMessage(signalId string, args Args) Message {
	return Message{
		MsgSignal,
		signalId,
		args,
	}
}

func MakeErrorMessage(msgType MsgType, id int64, error string) Message {
	return Message{
		MsgError,
		msgType,
		id,
		error,
	}
}
