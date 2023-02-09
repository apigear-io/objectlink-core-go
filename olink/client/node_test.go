package client

import (
	"encoding/json"
	"testing"

	"github.com/apigear-io/objectlink-core-go/olink/core"

	"github.com/stretchr/testify/assert"
)

func makeNodeAndSink(t *testing.T) (*Node, *MockSink, *core.MockDataWriter) {
	name := "demo.Counter"
	sink := &MockSink{objectId: name}
	registry := NewRegistry()
	writer := core.NewMockDataWriter()
	node := NewNode(registry)
	node.SetOutput(writer)
	return node, sink, writer
}

func TestClientGetSink(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.Registry().AddObjectSink(sink)
	sink2 := client.Registry().ObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
}

func TestClientRemoveSink(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.Registry().AddObjectSink(sink)
	client.Registry().RemoveObjectSink(sink.ObjectId())
	sink2 := client.Registry().ObjectSink(sink.ObjectId())
	assert.Nil(t, sink2, "sink should be nil")
}

func TestLinkNode(t *testing.T) {
	node, sink, _ := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	node.Registry().LinkClientNode(sink.ObjectId(), node)
	sink2 := node.Registry().ObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
}

func TestLinkRemote(t *testing.T) {
	node, sink, writer := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	// links and notifies remote
	node.LinkRemoteNode(sink.ObjectId())
	sink2 := node.Registry().ObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
}

func TestSetRemoteProperty(t *testing.T) {
	node, sink, writer := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	// links and notifies remote
	node.LinkRemoteNode(sink.ObjectId())
	sink2 := node.Registry().ObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
	// set property
	propertyId := core.MakeSymbolId(sink.ObjectId(), "prop")
	node.SetRemoteProperty(propertyId, "value")
	// writer should have one property message
	assert.Equal(t, 2, len(writer.Messages), "should have 2 message")

	propertyId2, value := writer.Messages[1].AsPropertyChange()
	assert.Equal(t, propertyId, propertyId2, "should be prop")
	assert.Equal(t, "value", value, "should be value")
}

func TestInvokeRemote(t *testing.T) {
	node, sink, writer := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	// links and notifies remote
	node.LinkRemoteNode(sink.ObjectId())
	sink2 := node.Registry().ObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
	// invoke remote
	methodId := core.MakeSymbolId(sink.ObjectId(), "method")
	node.InvokeRemote(methodId, core.Args{}, func(args InvokeReplyArg) {})
	// writer should have one invoke message
	assert.Equal(t, 2, len(writer.Messages), "should have 2 message")
	seq, res2, args := writer.Messages[1].AsInvoke()
	assert.Equal(t, int64(1), seq, "should be seq 1")
	assert.Equal(t, methodId, res2, "should be method")
	assert.Equal(t, core.Args{}, args, "should be args")
}

func TestUnlinkRemoteNode(t *testing.T) {
	node, sink, writer := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	// links and notifies remote
	node.LinkRemoteNode(sink.ObjectId())
	sink2 := node.Registry().ObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
	// unlink remote
	node.UnlinkRemoteNode(sink.ObjectId())
	// writer should have one unlink message
	assert.Equal(t, 2, len(writer.Messages), "should have 2 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[1].AsUnlink(), "should be unlink message")
}

func TestHandleInit(t *testing.T) {
	node, sink, _ := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	node.Registry().LinkClientNode(sink.ObjectId(), node)
	msg := core.MakeInitMessage(sink.ObjectId(), core.KWArgs{})
	data, err := json.Marshal(msg)
	assert.Nil(t, err, "should be nil")
	node.Write(data)
	assert.Equal(t, 1, len(sink.events), "should have 1 event")
	assert.Equal(t, msg, sink.events[0], "should be init event")
}

func TestHandlePropertyChange(t *testing.T) {
	node, sink, _ := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	node.Registry().LinkClientNode(sink.ObjectId(), node)
	propertyId := core.MakeSymbolId(sink.ObjectId(), "prop")
	msg := core.MakePropertyChangeMessage(propertyId, "value")
	data, err := json.Marshal(msg)
	assert.Nil(t, err, "should be nil")
	node.Write(data)
	assert.Equal(t, 1, len(sink.events), "should have 1 event")
	assert.Equal(t, msg, sink.events[0], "should be property event")
}

func TestHandleMsgInvokeReply(t *testing.T) {
	node, sink, _ := makeNodeAndSink(t)
	node.Registry().AddObjectSink(sink)
	node.Registry().LinkClientNode(sink.ObjectId(), node)
	isCalled := false
	methodId := core.MakeSymbolId(sink.ObjectId(), "hello")
	node.InvokeRemote(methodId, core.Args{}, func(args InvokeReplyArg) {
		isCalled = true
	})
	msg := core.MakeInvokeReplyMessage(1, methodId, "value")
	data, err := json.Marshal(msg)
	assert.Nil(t, err, "should be nil")
	node.Write(data)
	assert.True(t, isCalled, "should be called")
}
