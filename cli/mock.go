package cli

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

type MockSink struct {
	node     *client.Node
	objectId string
	props    core.KWArgs
}

func (s *MockSink) ObjectId() string {
	return s.objectId
}

func (s *MockSink) HandleSignal(signalId string, args core.Args) {
	fmt.Printf("%s: signal: %s %v\n", s.ObjectId(), signalId, args)
}

func (s *MockSink) HandlePropertyChange(propertyId string, value core.Any) {
	fmt.Printf("%s: property change: %s %v\n", s.ObjectId(), propertyId, value)
}

func (s *MockSink) HandleInit(objectId string, props core.KWArgs, node *client.Node) {
	data, err := json.MarshalIndent(props, "", "  ")
	if err != nil {
		log.Printf("error marshalling value: %v", err)
		return
	}
	fmt.Printf("on init %s\n", objectId)
	fmt.Println(string(data))
	if objectId != s.ObjectId() {
		return
	}
	s.props = props
	s.node = node
}

func (s *MockSink) HandleRelease() {
	fmt.Printf("%s: on release\n", s.ObjectId())
	s.node = nil
}
