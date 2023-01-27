package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/remote"
	"github.com/apigear-io/objectlink-core-go/olink/ws"
)

var addr = flag.String("addr", ":8080", "http service address")

type INotifier interface {
	NotifyPropertyChanged(name string, value core.Any)
	NotifySignal(name string, args core.Args)
}

type ICounter interface {
	INotifier
	GetCount() int64
	SetCount(count int64)
	Increment(step int64)
	Decrement(step int64)
}

type counterImpl struct {
	INotifier
	count int64
}

var _ ICounter = (*counterImpl)(nil)
var _ INotifier = (*counterImpl)(nil)

func NewCounter(notifier INotifier) ICounter {
	return &counterImpl{
		INotifier: notifier,
	}
}

func (impl *counterImpl) GetCount() int64 {
	log.Info().Msgf("impl: get count: %d", impl.count)
	return impl.count
}

func (impl *counterImpl) SetCount(count int64) {
	log.Info().Msgf("impl: set count: %d", count)
	impl.count = count
	impl.NotifyPropertyChanged("count", impl.count)
}

func (impl *counterImpl) Increment(step int64) {
	log.Info().Msgf("impl: increment: %d", step)
	impl.count += step
	impl.NotifyPropertyChanged("count", impl.count)
}

func (impl *counterImpl) Decrement(step int64) {
	log.Info().Msgf("impl: decrement: %d", step)
	impl.count -= step
	impl.NotifyPropertyChanged("count", impl.count)
}

type CounterSource struct {
	node *remote.Node
	impl ICounter
}

var _ remote.IObjectSource = (*CounterSource)(nil)
var _ INotifier = (*CounterSource)(nil)

func NewCounterSource() *CounterSource {
	return &CounterSource{}
}

func (s *CounterSource) SetImplementation(impl ICounter) {
	s.impl = impl
}

func (s *CounterSource) ObjectId() string {
	return "demo.Counter"
}

func (s *CounterSource) Invoke(methodId string, args core.Args) (core.Any, error) {
	log.Info().Msgf("source: invoke: %s %v", methodId, args)
	if s.impl == nil {
		return nil, fmt.Errorf("no implementation")
	}
	name := core.SymbolIdToMember(methodId)
	switch name {
	case "increment":
		s.impl.Increment(args[0].(int64))
		return nil, nil
	case "decrement":
		s.impl.Decrement(args[0].(int64))
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown method: %s", name)
	}
}
func (s *CounterSource) SetProperty(propertyId string, value core.Any) error {
	log.Info().Msgf("source: set property %s %v", propertyId, value)
	if s.impl == nil {
		return fmt.Errorf("no implementation")
	}
	name := core.SymbolIdToMember(propertyId)
	switch name {
	case "count":
		s.impl.SetCount(value.(int64))
	default:
		return fmt.Errorf("unknown property: %s", name)
	}
	return nil
}
func (s *CounterSource) Linked(objectId string, node *remote.Node) error {
	log.Info().Msgf("source: linked %s %v", objectId, node)
	if objectId != s.ObjectId() {
		return fmt.Errorf("unexpected object id: %s", objectId)
	}
	if s.node != nil {
		return fmt.Errorf("already linked")
	}
	s.node = node
	return nil
}

func (s *CounterSource) CollectProperties() (core.KWArgs, error) {
	log.Info().Msgf("source: collect properties")
	if s.impl == nil {
		return nil, fmt.Errorf("no implementation")
	}
	return core.KWArgs{
		"count": s.impl.GetCount(),
	}, nil
}

func (s *CounterSource) NotifyPropertyChanged(name string, value core.Any) {
	propertyId := core.MakeSymbolId(s.ObjectId(), name)
	s.node.NotifyPropertyChange(propertyId, value)
}

func (s *CounterSource) NotifySignal(name string, args core.Args) {
	signalId := core.MakeSymbolId(s.ObjectId(), name)
	s.node.NotifySignal(signalId, args)
}

func main() {
	flag.Parse()
	registry := remote.NewRegistry()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	hub := ws.NewHub(ctx, registry)
	{
		source := NewCounterSource()
		registry.AddObjectSource(source)
		impl := NewCounter(source)
		source.SetImplementation(impl)
	}

	http.HandleFunc("/ws", hub.ServeHTTP)
	log.Info().Msgf("web socket server listening on %s/ws", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start web socket server")
	}
}
