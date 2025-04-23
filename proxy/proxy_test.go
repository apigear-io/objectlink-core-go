package proxy

import (
	"fmt"
	"io"
	"testing"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/remote"
	"github.com/stretchr/testify/assert"
)

type DemoSource struct {
	id         string
	value      int
	codec      Codec
	node       *remote.Node
	properties core.KWArgs
	signals    map[string]core.Args
	invokes    map[string]core.Args
}

var _ remote.IObjectSource = (*DemoSource)(nil)

func NewDemoSource(id string) *DemoSource {
	return &DemoSource{
		id:    id,
		codec: NewCodec("json"),
		properties: core.KWArgs{
			"value": 0,
		},
		signals: make(map[string]core.Args),
		invokes: make(map[string]core.Args),
	}
}

func (s *DemoSource) ObjectId() string {
	return s.id
}

func (s *DemoSource) NotifyPropertyChanges(name string, value core.Any) {
	if s.node == nil {
		return
	}
	s.properties[name] = value
	propertyId := core.MakeSymbolId(s.id, name)
	s.node.NotifyPropertyChange(propertyId, value)
}

func (s *DemoSource) NotifySignal(name string, args core.Args) {
	if s.node == nil {
		return
	}
	s.signals[name] = args
	signalId := core.MakeSymbolId(s.id, name)
	s.node.NotifySignal(signalId, args)
}

func (s *DemoSource) Invoke(methodId string, args core.Args) (core.Any, error) {
	s.invokes[methodId] = args
	switch methodId {
	case "inc":
		step := 1

		if len(args) > 0 {
			step = args[0].(int)
		}
		s.value += step
		s.NotifyPropertyChanges("value", s.value)
		return s.value, nil
	case "dec":
		step := 1
		if len(args) > 0 {
			step = args[0].(int)
		}
		s.value -= step
		s.NotifyPropertyChanges("value", s.value)
		return s.value, nil
	case "reset":
		s.value = 0
		s.NotifyPropertyChanges("value", s.value)
		s.NotifySignal("reset", core.Args{})
		return s.value, nil
	}

	return nil, fmt.Errorf("invalid method: %s", methodId)
}

func (s *DemoSource) SetProperty(propertyId string, value core.Any) error {
	member := core.SymbolIdToMember(propertyId)
	s.properties[member] = value
	switch member {
	case "value":
		s.value = value.(int)
	}
	return nil
}

func (s *DemoSource) Linked(objectId string, node *remote.Node) error {
	if objectId != s.id {
		return fmt.Errorf("invalid object id: %s", objectId)
	}
	s.node = node
	return nil
}
func (s *DemoSource) CollectProperties() (core.KWArgs, error) {
	return s.properties, nil
}

type DemoSink struct {
	id         string
	value      int
	node       *client.Node
	signals    map[string]core.Args
	properties core.KWArgs
}

var _ client.IObjectSink = (*DemoSink)(nil)

func (s *DemoSink) ObjectId() string {
	return s.id
}

func NewDemoSink(id string) *DemoSink {
	return &DemoSink{
		id:         id,
		signals:    make(map[string]core.Args),
		properties: core.KWArgs{},
	}
}

func (s *DemoSink) Inc(step int) {
	if s.node == nil {
		return
	}
	methodId := core.MakeSymbolId(s.id, "inc")
	args := core.Args{step}
	s.node.InvokeRemote(methodId, args, func(arg client.InvokeReplyArg) {
		log.Info().Msgf("Inc: %v", arg)
	})
	return
}

func (s *DemoSink) Dec(step int) {
	if s.node == nil {
		return
	}
	methodId := core.MakeSymbolId(s.id, "dec")
	args := core.Args{step}
	s.node.InvokeRemote(methodId, args, func(arg client.InvokeReplyArg) {
		log.Info().Msgf("Dec: %v", arg)
	})
}

func (s *DemoSink) Reset() {
	if s.node == nil {
		return
	}
	methodId := core.MakeSymbolId(s.id, "reset")
	args := core.Args{}
	s.node.InvokeRemote(methodId, args, func(arg client.InvokeReplyArg) {
		log.Info().Msgf("Reset: %v", arg)
	})
}

func (s *DemoSink) HandleSignal(signalId string, args core.Args) {
	log.Info().Msgf("OnSignal: %s", signalId)
	s.signals[signalId] = args
}
func (s *DemoSink) HandlePropertyChange(propertyId string, value core.Any) {
	log.Info().Msgf("OnPropertyChange: %s", propertyId)
	s.properties[propertyId] = value
}
func (s *DemoSink) HandleInit(objectId string, props core.KWArgs, node *client.Node) {
	log.Info().Msgf("OnInit: %s", objectId)
	s.properties = props
	s.node = node
}
func (s *DemoSink) HandleRelease() {
	log.Info().Msgf("OnRelease")
	s.node = nil
}

func writeTo(t *testing.T, w io.Writer, msg core.Message) {
	conv := core.NewConverter(core.FormatJson)
	data, err := conv.ToData(msg)
	assert.NoError(t, err)
	w.Write(data)
}

// func TestLinkMessages(t *testing.T) {
// 	conv := core.NewConverter(core.FormatJson)
// 	state := core.KWArgs{"value": 10}
// 	unlink := core.MakeUnlinkMessage("demo")
// 	init := core.MakeInitMessage("demo", state)
// 	proxy := NewProxy()
// 	rn := proxy.CreateRemoteNode()
// 	cn := proxy.CreateClientNode()
// 	writeTo(t, rn, core.MakeLinkMessage("demo"))
// 	proxy.HandleMessage(link)
// 	proxy.OnMessage(func(msg core.Message) {
// 		t.Errorf("unexpected message: %v", msg)
// 	})
// 	source := proxy.Source()
// 	sink := proxy.Sink()
// 	if source == nil {
// 		t.Errorf("source is nil")
// 	}
// 	if sink == nil {
// 		t.Errorf("sink is nil")
// 	}
// }
