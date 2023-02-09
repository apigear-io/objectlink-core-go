package client

import "github.com/apigear-io/objectlink-core-go/olink/core"

// MockSink is an IObjectSink implementation
type MockSink struct {
	events   []core.Message
	objectId string
}

func NewMockSink(objectId string) *MockSink {
	return &MockSink{
		events:   make([]core.Message, 0),
		objectId: objectId,
	}
}

func (m *MockSink) ObjectId() string {
	return m.objectId
}

func (m *MockSink) OnSignal(signalId string, args core.Args) {
	m.events = append(m.events, core.MakeSignalMessage(signalId, args))
}

func (m *MockSink) OnPropertyChange(propertyId string, value core.Any) {
	m.events = append(m.events, core.MakePropertyChangeMessage(propertyId, value))
}

func (m *MockSink) OnInit(objectId string, props core.KWArgs, node *Node) {
	m.events = append(m.events, core.MakeInitMessage(objectId, props))
}

func (m *MockSink) OnRelease() {}
