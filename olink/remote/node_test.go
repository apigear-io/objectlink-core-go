package remote

import (
	"testing"

	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	r := NewRegistry()
	n := NewNode(r)
	assert.NotNil(t, n)
	assert.NotNil(t, n.Registry())
}

func TestNodeSetOutput(t *testing.T) {
	r := NewRegistry()
	n := NewNode(r)
	assert.NotNil(t, n)
	assert.NotNil(t, n.Registry())
	wc := NewMockWriteCloser()
	n.SetOutput(wc)
	msg := core.MakeLinkMessage("test")
	n.SendMessage(msg)
	assert.Equal(t, 1, len(wc.Messages))
	act, err := n.conv.FromData(wc.Messages[0])
	assert.Nil(t, err)
	assert.Equal(t, msg.AsLink(), act.AsLink())
}

func TestNodeSendMessage(t *testing.T) {
	r := NewRegistry()
	n := NewNode(r)
	assert.NotNil(t, n)
	assert.NotNil(t, n.Registry)
	wc := NewMockWriteCloser()
	n.SetOutput(wc)
	msg := core.MakeLinkMessage("test")
	n.SendMessage(msg)
	assert.Equal(t, 1, len(wc.Messages))
	act, err := n.conv.FromData(wc.Messages[0])
	assert.Nil(t, err)
	assert.Equal(t, msg.AsLink(), act.AsLink())
}

// TestRemoveNode is a test for Node.RemoveNode()
func TestRemoveNode(t *testing.T) {
	r := NewRegistry()
	n := NewNode(r)
	s := NewMockSource("test")
	r.AddObjectSource(s)
	r.LinkRemoteNode(s.ObjectId(), n)
	assert.True(t, r.IsRegistered(s.ObjectId()))
	assert.Equal(t, 1, len(r.GetRemoteNodes(s.ObjectId())))
	n.RemoveNode()
	assert.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
	n.RemoveNode()
	assert.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
}
