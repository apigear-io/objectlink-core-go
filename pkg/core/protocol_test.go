package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var name = "demo.Calc"
var props = Props{"count": 1}
var value = 1
var args = Args{1, 3}
var lastValue = 0
var requestId = 1
var msgType = INVOKE
var errorMessage = "failed"

func TestLinkMessage(t *testing.T) {
	msg := NewLinkMessage(name)
	assert.Equal(t, Args{LINK, name}, msg)
}

func TestUnlinkMessage(t *testing.T) {
	msg := NewUnlinkMessage(name)
	assert.Equal(t, Args{UNLINK, name}, msg)
}

func TestInitMessage(t *testing.T) {
	msg := NewInitMessage(name, props)
	assert.Equal(t, Args{INIT, name, props}, msg)
}

func TestSetProperty(t *testing.T) {
	msg := NewSetPropertyMessage(name, value)
	assert.Equal(t, Args{SET_PROPERTY, name, value}, msg)
}

func TestPropertyChange(t *testing.T) {
	msg := NewPropertyChangeMessage(name, value)
	assert.Equal(t, Args{PROPERTY_CHANGE, name, value}, msg)
}

func TestInvoke(t *testing.T) {
	msg := NewInvokeMessage(requestId, name, args)
	assert.Equal(t, Args{INVOKE, requestId, name, args}, msg)
}

func TestInvokeReply(t *testing.T) {
	msg := NewInvokeReplyMessage(requestId, name, value)
	assert.Equal(t, Args{INVOKE_REPLY, requestId, name, value}, msg)
}

func TestSignal(t *testing.T) {
	msg := NewSignalMessage(name, args)
	assert.Equal(t, Args{SIGNAL, name, args}, msg)
}

func TestError(t *testing.T) {
	msg := NewErrorMessage(msgType, requestId, errorMessage)
	assert.Equal(t, Args{ERROR, msgType, requestId, errorMessage}, msg)
}
