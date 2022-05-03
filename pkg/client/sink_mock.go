package client

import "olink/pkg/core"

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
