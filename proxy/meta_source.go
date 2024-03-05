package proxy

import (
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/remote"
)

type MetaSource struct {
	id         string
	codec      Codec
	node       *remote.Node
	Properties core.KWArgs
	OnInvoke   map[string]func(args core.Args) core.Any
}

var _ remote.IObjectSource = (*MetaSource)(nil)

func (m *MetaSource) ObjectId() string {
	return m.id
}
func (m *MetaSource) Invoke(methodId string, args core.Args) (core.Any, error) {
	result := m.OnInvoke[methodId](args)
	if m.node == nil {
		return nil, nil
	}
	return result, nil
}
func (m *MetaSource) SetProperty(propertyId string, value core.Any) error {
	name := core.SymbolIdToMember(propertyId)
	m.Properties[name] = value
	return nil
}
func (m *MetaSource) Linked(objectId string, node *remote.Node) error {
	m.node = node
	return nil
}

func (m *MetaSource) CollectProperties() (core.KWArgs, error) {
	return m.Properties, nil
}
