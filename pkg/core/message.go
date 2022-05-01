package core

type MsgType int

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
type Props map[string]any
type Any any

type Message []any

func (m Message) Type() MsgType {
	return AsMsgType(m[0])
}

// AsInit returns the name and props of the init message
func (m Message) AsInit() (string, Props) {
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
func (m Message) AsSetProperty() (Resource, Any) {
	return AsResource(m[1]), AsAny(m[2])
}

// AsPropertyChange returns the name and value of the property change
func (m Message) AsPropertyChange() (Resource, Any) {
	return AsResource(m[1]), AsAny(m[2])
}

// AsInvoke returns the id, name and args of the invoke message
func (m Message) AsInvoke() (int, Resource, Args) {
	return AsInt(m[1]), AsResource(m[2]), AsArgs(m[3])
}

// AsInvokeReply returns the id and result of the invoke reply
func (m Message) AsInvokeReply() (int, Resource, Any) {
	return AsInt(m[1]), AsResource(m[2]), AsAny(m[3])
}

// AsSignal returns the name and args of the signal message
func (m Message) AsSignal() (Resource, Args) {
	return AsResource(m[1]), AsArgs(m[2])
}

func (m Message) AsError() (MsgType, int, string) {
	return AsMsgType(m[0]), AsInt(m[1]), AsString(m[2])
}

func CreateLinkMessage(objectId string) Message {
	return Message{
		MsgLink,
		objectId,
	}
}

func CreateInitMessage(objectId string, props Props) Message {
	return Message{
		MsgInit,
		objectId,
		props,
	}
}

func CreateUnlinkMessage(objectId string) Message {
	return Message{
		MsgUnlink,
		objectId,
	}
}

func CreateSetPropertyMessage(res Resource, value Any) Message {
	return Message{
		MsgSetProperty,
		res,
		value,
	}
}

func CreatePropertyChangeMessage(res Resource, value Any) Message {
	return Message{
		MsgPropertyChange,
		res,
		value,
	}
}

func CreateInvokeMessage(id int, res Resource, args Args) Message {
	return Message{
		MsgInvoke,
		id,
		res,
		args,
	}
}

func CreateInvokeReplyMessage(id int, res Resource, value Any) Message {
	return Message{
		MsgInvokeReply,
		id,
		res,
		value,
	}
}

func CreateSignalMessage(res Resource, args Args) Message {
	return Message{
		MsgSignal,
		res,
		args,
	}
}

func CreateErrorMessage(msgType MsgType, id int, error string) Message {
	return Message{
		MsgError,
		msgType,
		id,
		error,
	}
}
