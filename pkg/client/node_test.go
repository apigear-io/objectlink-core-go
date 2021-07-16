package client

import (
	"olink/pkg/core"
	"testing"
)

// MockSink is an IObjectSink implementation
type MockSink struct {
	events []core.Message
	name   string
}

func (m *MockSink) ObjectName() string {
	return m.name
}

func (m *MockSink) OnSignal(name string, args core.Args) {
	m.events = append(m.events, core.NewSignalMessage(name, args))
}

func (m *MockSink) OnPropertyChange(name string, value core.Any) {
	m.events = append(m.events, core.NewPropertyChangeMessage(name, value))
}

func (m *MockSink) OnInit(name string, props core.Props, node *ClientNode) {
	m.events = append(m.events, core.NewInitMessage(name, props))
}

func (m *MockSink) OnRelease() {}

var name = "demo.Counter"
var sink = &MockSink{name: name}
var client = NewClientNode()
var r = GetRegistry()

func TestAddSink(t *testing.T) {
	client.AddObjectSink(sink)
}
