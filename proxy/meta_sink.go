package proxy

import (
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

type MetaSink struct {
	id         string
	codec      Codec
	node       *client.Node
	Properties core.KWArgs
	Change     func(name string, value core.Any)
	Signal     func(name string, args core.Args)
}

func NewMetaSink(id string) *MetaSink {
	return &MetaSink{
		id:         id,
		codec:      NewCodec("json"),
		Properties: make(core.KWArgs),
	}
}

var _ client.IObjectSink = (*MetaSink)(nil)

func (s *MetaSink) ObjectId() string {
	return s.id
}

func (s *MetaSink) OnSignal(signalId string, args core.Args) {
	if s.Signal != nil {
		s.Signal(signalId, args)
	}
}
func (s *MetaSink) OnPropertyChange(propertyId string, value core.Any) {
	name := core.SymbolIdToMember(propertyId)
	s.Properties[name] = value
	if s.Change != nil {
		s.Change(name, value)
	}
}
func (s *MetaSink) OnInit(objectId string, props core.KWArgs, node *client.Node) {
	s.Properties = props
	s.node = node

}
func (s *MetaSink) OnRelease() {
	s.node = nil
}

func (s *MetaSink) Invoke(methodId string, args core.Args, cb func(v core.Any)) {
	if s.node == nil {
		return
	}
	fn := func(arg client.InvokeReplyArg) {
		if cb != nil {
			cb(arg.Value)
		}
	}
	s.node.InvokeRemote(methodId, args, fn)
}
