package client

import (
	"olink/pkg/core"
	"testing"
)

type CounterSink struct {
	events []core.Message
	Count  int
}

func (s *CounterSink) ObjectId() string {
	return "demo.Counter"
}

func (s *CounterSink) OnSignal(signalId string, args core.Args) {
	s.events = append(s.events, core.MakeSignalMessage(signalId, args))
}

func (s *CounterSink) OnPropertyChange(propertyId string, value core.Any) {
	s.events = append(s.events, core.MakePropertyChangeMessage(propertyId, value))
	if core.ToMember(propertyId) == "count" {
		s.Count = value.(int)
	}
}

func (s *CounterSink) OnInit(objectId string, props core.KWArgs, node *Node) {
	s.events = append(s.events, core.MakeInitMessage(objectId, props))
	_, ok := props["count"]
	if ok {
		s.Count = props["count"].(int)
	}
}

func (s *CounterSink) OnRelease() {}

func TestCounterSink(t *testing.T) {
	sink := &CounterSink{}
	writer := core.NewMockDataWriter()
	registry := NewRegistry()
	node := NewNode(registry)
	node.SetOutput(writer)

	// link node to sink object id
	// registry.LinkClientNode(sink.ObjectId(), node)
	registry.LinkClientNode(sink.ObjectId(), node)
	// register sink using objectId
	// registry.AddObjectSink(sink)
	registry.AddObjectSink(sink)
	// subscribe to remote source events
	node.LinkRemoteNode(sink.ObjectId())
	res := core.MakeIdentifier(sink.ObjectId(), "count")
	node.SetRemoteProperty(res, 0)
}
