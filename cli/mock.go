package cli

import (
	"fmt"
	"olink/pkg/client"
	"olink/pkg/core"
)

type MockSink struct {
	node     *client.Node
	objectId string
	props    core.KWArgs
}

func (s *MockSink) ObjectId() string {
	return s.objectId
}

func (s *MockSink) OnSignal(signalId string, args core.Args) {
	fmt.Printf("%s: signal: %s %v\n", s.ObjectId(), signalId, args)
}

func (s *MockSink) OnPropertyChange(propertyId string, value core.Any) {
	fmt.Printf("%s: property change: %s %v\n", s.ObjectId(), propertyId, value)
}

func (s *MockSink) OnInit(objectId string, props core.KWArgs, node *client.Node) {
	fmt.Printf("%s: on init %s %#v\n", s.ObjectId(), objectId, props)
	if objectId != s.ObjectId() {
		return
	}
	s.props = props
	s.node = node
}

func (s *MockSink) OnRelease() {
	fmt.Printf("%s: on release\n", s.ObjectId())
	s.node = nil
}
