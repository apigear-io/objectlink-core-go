package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	ObjectId     string
	Resource     Resource
	Props        Props
	Value        interface{}
	LastValue    interface{}
	Args         Args
	RequestId    int
	ErrorMessage string
	MsgType      MsgType
}

var data = Data{
	ObjectId: "demo.calc",
	Resource: "demo.calc/add",
	Props: map[string]interface{}{
		"count": 1,
	},
	Value:        1,
	LastValue:    0,
	Args:         []any{1, 3},
	RequestId:    1,
	ErrorMessage: "error",
}

// var name = "demo.Calc"
// var props = Props{"count": 1}
// var value = 1
// var args = Args{1, 3}
// var lastValue = 0
// var requestId = 1
// var msgType = MsgInvoke
// var errorMessage = "failed"

func TestLinkMessage(t *testing.T) {
	msg := CreateLinkMessage(data.ObjectId)
	assert.Equal(t, Message{MsgLink, data.ObjectId}, msg)
}

func TestUnlinkMessage(t *testing.T) {
	msg := CreateUnlinkMessage(data.ObjectId)
	assert.Equal(t, Message{MsgUnlink, data.ObjectId}, msg)
}

func TestInitMessage(t *testing.T) {
	msg := CreateInitMessage(data.ObjectId, data.Props)
	assert.Equal(t, Message{MsgInit, data.ObjectId, data.Props}, msg)
}

func TestSetProperty(t *testing.T) {
	msg := CreateSetPropertyMessage(data.Resource, data.Value)
	assert.Equal(t, Message{MsgSetProperty, data.Resource, data.Value}, msg)
}

func TestPropertyChange(t *testing.T) {
	msg := CreatePropertyChangeMessage(data.Resource, data.Value)
	assert.Equal(t, Message{MsgPropertyChange, data.Resource, data.Value}, msg)
}

func TestInvoke(t *testing.T) {
	msg := CreateInvokeMessage(data.RequestId, data.Resource, data.Args)
	assert.Equal(t, Message{MsgInvoke, data.RequestId, data.Resource, data.Args}, msg)
}

func TestInvokeReply(t *testing.T) {
	msg := CreateInvokeReplyMessage(data.RequestId, data.Resource, data.Value)
	assert.Equal(t, Message{MsgInvokeReply, data.RequestId, data.Resource, data.Value}, msg)
}

func TestSignal(t *testing.T) {
	msg := CreateSignalMessage(data.Resource, data.Args)
	assert.Equal(t, Message{MsgSignal, data.Resource, data.Args}, msg)
}

func TestError(t *testing.T) {
	msg := CreateErrorMessage(data.MsgType, data.RequestId, data.ErrorMessage)
	assert.Equal(t, Message{MsgError, data.MsgType, data.RequestId, data.ErrorMessage}, msg)
}
