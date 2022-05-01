package client

import (
	"encoding/json"
	"olink/pkg/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDataWriter struct {
	Messages  []core.Message
	converter *core.MessageConverter
}

func NewMockDataWriter() *MockDataWriter {
	return &MockDataWriter{
		converter: &core.MessageConverter{Format: core.FormatJson},
	}
}

func (w *MockDataWriter) WriteData(data []byte) error {
	msg, err := w.converter.FromData(data)
	if err != nil {
		return err
	}
	w.Messages = append(w.Messages, msg)
	return nil
}

// MockSink is an IObjectSink implementation
type MockSink struct {
	events   []core.Message
	objectId string
}

func (m *MockSink) ObjectId() string {
	return m.objectId
}

func (m *MockSink) OnSignal(res core.Resource, args core.Args) {
	m.events = append(m.events, core.CreateSignalMessage(res, args))
}

func (m *MockSink) OnPropertyChange(res core.Resource, value core.Any) {
	m.events = append(m.events, core.CreatePropertyChangeMessage(res, value))
}

func (m *MockSink) OnInit(objectId string, props core.Props, node *Node) {
	m.events = append(m.events, core.CreateInitMessage(objectId, props))
}

func (m *MockSink) OnRelease() {}

func makeNodeAndSink(t *testing.T) (*Node, *MockSink, *MockDataWriter) {
	name := "demo.Counter"
	sink := &MockSink{objectId: name}
	registry := NewRegistry()
	writer := NewMockDataWriter()
	client := NewNode(registry, writer)
	return client, sink, writer
}

func TestClientGetSink(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
}

func TestClientRemoveSink(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	client.RemoveObjectSink(sink)
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Nil(t, sink2, "sink should be nil")
}

func TestLinkNode(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	client.LinkNode(sink.ObjectId())
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
}

func TestLinkRemote(t *testing.T) {
	client, sink, writer := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	// links and notifies remote
	client.LinkRemoteNode(sink.ObjectId())
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
}

func TestSetRemoteProperty(t *testing.T) {
	client, sink, writer := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	// links and notifies remote
	client.LinkRemoteNode(sink.ObjectId())
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
	// set property
	res := core.CreateResource(sink.ObjectId(), "prop")
	client.SetRemoteProperty(res, "value")
	// writer should have one property message
	assert.Equal(t, 2, len(writer.Messages), "should have 2 message")

	res2, value := writer.Messages[1].AsPropertyChange()
	assert.Equal(t, res, res2, "should be prop")
	assert.Equal(t, "value", value, "should be value")
}

func TestInvokeRemote(t *testing.T) {
	client, sink, writer := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	// links and notifies remote
	client.LinkRemoteNode(sink.ObjectId())
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
	// invoke remote
	res := core.CreateResource(sink.ObjectId(), "method")
	client.InvokeRemote(res, core.Args{}, func(args InvokeReplyArg) {})
	// writer should have one invoke message
	assert.Equal(t, 2, len(writer.Messages), "should have 2 message")
	seq, res2, args := writer.Messages[1].AsInvoke()
	assert.Equal(t, 1, seq, "should be seq 1")
	assert.Equal(t, res, res2, "should be method")
	assert.Equal(t, core.Args{}, args, "should be args")
}

func TestUnlinkRemoteNode(t *testing.T) {
	client, sink, writer := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	// links and notifies remote
	client.LinkRemoteNode(sink.ObjectId())
	sink2 := client.GetObjectSink(sink.ObjectId())
	assert.Equal(t, sink, sink2, "sink should be the same")
	// writer should have one link message
	assert.Equal(t, 1, len(writer.Messages), "should have 1 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[0].AsLink(), "should be init message")
	// unlink remote
	client.UnlinkRemoteNode(sink.ObjectId())
	// writer should have one unlink message
	assert.Equal(t, 2, len(writer.Messages), "should have 2 message")
	assert.Equal(t, sink.ObjectId(), writer.Messages[1].AsUnlink(), "should be unlink message")
}

func TestHandleInit(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	client.LinkNode(sink.ObjectId())
	msg := core.CreateInitMessage(sink.ObjectId(), core.Props{})
	data, err := json.Marshal(msg)
	assert.Nil(t, err, "should be nil")
	client.HandleMessage(data)
	assert.Equal(t, 1, len(sink.events), "should have 1 event")
	assert.Equal(t, msg, sink.events[0], "should be init event")
}

func TestHandlePropertyChange(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	client.LinkNode(sink.ObjectId())
	res := core.CreateResource(sink.ObjectId(), "prop")
	msg := core.CreatePropertyChangeMessage(res, "value")
	data, err := json.Marshal(msg)
	assert.Nil(t, err, "should be nil")
	client.HandleMessage(data)
	assert.Equal(t, 1, len(sink.events), "should have 1 event")
	assert.Equal(t, msg, sink.events[0], "should be property event")
}

func TestHandleMsgInvokeReply(t *testing.T) {
	client, sink, _ := makeNodeAndSink(t)
	client.AddObjectSink(sink)
	client.LinkNode(sink.ObjectId())
	isCalled := false
	res := core.CreateResource(sink.ObjectId(), "hello")
	client.InvokeRemote(res, core.Args{}, func(args InvokeReplyArg) {
		isCalled = true
	})
	msg := core.CreateInvokeReplyMessage(1, res, "value")
	data, err := json.Marshal(msg)
	assert.Nil(t, err, "should be nil")
	client.HandleMessage(data)
	assert.True(t, isCalled, "should be called")
}
