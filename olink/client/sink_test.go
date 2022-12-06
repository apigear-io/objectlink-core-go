package client

import (
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"testing"
)

type CounterSink struct {
	events []core.Message
	Count  int
}

func (s *CounterSink) ObjectId() string {
	return "demo.Counter"
}

func (s *CounterSink) OnSignal(res core.Resource, args core.Args) {
	s.events = append(s.events, core.CreateSignalMessage(res, args))
}

func (s *CounterSink) OnPropertyChange(res core.Resource, value core.Any) {
	s.events = append(s.events, core.CreatePropertyChangeMessage(res, value))
	if res.Member() == "Count" {
		s.Count = value.(int)
	}
}

func (s *CounterSink) OnInit(objectId string, props core.Props, node *Node) {
	s.events = append(s.events, core.CreateInitMessage(objectId, props))
	_, ok := props["count"]
	if ok {
		s.Count = props["count"].(int)
	}
}

func (s *CounterSink) OnRelease() {}

func TestCounterSink(t *testing.T) {
	sink := &CounterSink{}
	writer := core.NewMockWriter()
	registry := NewRegistry()
	node := NewNode(registry, writer)

	// link node to sink object id
	// registry.LinkClientNode(sink.ObjectId(), node)
	registry.LinkClientNode(sink.ObjectId(), node)
	// register sink using objectId
	// registry.AddObjectSink(sink)
	registry.AddObjectSink(sink)
	// subscribe to remote source events
	node.LinkRemoteNode(sink.ObjectId())
	res := core.CreateResource(sink.ObjectId(), "count")
	node.SetRemoteProperty(res, 0)
}
