package remote

import "github.com/apigear-io/objectlink-core-go/olink/core"

type MockSource struct {
	Id                       string
	Messages                 core.Message
	InvokeHandler            func(methodId string, args core.Args) (core.Any, error)
	SetPropertyHandler       func(propertyId string, value core.Any) error
	LinkedHandler            func(objectId string, node *Node) error
	CollectPropertiesHandler func() (core.KWArgs, error)
}

func NewMockSource(id string) *MockSource {
	return &MockSource{
		Id: id,
	}
}

func (s *MockSource) ObjectId() string {
	return s.Id
}

func (s *MockSource) Invoke(methodId string, args core.Args) (core.Any, error) {
	if s.InvokeHandler != nil {
		return s.InvokeHandler(methodId, args)
	}
	return nil, nil
}

func (s *MockSource) SetProperty(propertyId string, value core.Any) error {
	if s.SetPropertyHandler != nil {
		return s.SetPropertyHandler(propertyId, value)
	}
	return nil
}

func (s *MockSource) CollectProperties() (core.KWArgs, error) {
	if s.CollectPropertiesHandler != nil {
		return s.CollectPropertiesHandler()
	}
	return nil, nil
}

func (s *MockSource) Linked(objectId string, node *Node) error {
	if s.LinkedHandler != nil {
		return s.LinkedHandler(objectId, node)
	}
	return nil
}
