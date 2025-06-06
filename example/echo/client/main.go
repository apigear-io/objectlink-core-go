package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/ws"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/client"
)

var addr = flag.String("addr", "ws://127.0.0.1:8080/ws", "http ws service address")

// type CounterIncrementRequest struct {
// 	Step int `json:"step"`
// }

// type CounterIncrementReply struct {
// }

// type CounterDecrementRequest struct {
// 	Step int `json:"step"`
// }

// type CounterDecrementReply struct {
// }

// type CounterProperties struct {
// 	Count *int `json:"count,omitempty"`
// }

type CounterSink struct {
	count int64
	node  *client.Node
}

var _ client.IObjectSink = (*CounterSink)(nil)

func NewSink(node *client.Node) *CounterSink {
	return &CounterSink{
		node: node,
	}
}

func (s *CounterSink) ObjectId() string {
	return "demo.Counter"
}

func (s *CounterSink) SetCount(count int64) {
	if s.node != nil {
		propertyId := core.MakeSymbolId(s.ObjectId(), "count")
		s.node.SetRemoteProperty(propertyId, s.count)
	}
}

func (s *CounterSink) Increment(step int64) {
	log.Info().Msgf("sink: increment %s: %d", s.ObjectId(), step)
	if s.node == nil {
		log.Info().Msgf("no node")
		return
	}
	s.node.InvokeRemote(core.MakeSymbolId(s.ObjectId(), "increment"), core.Args{step}, nil)
}

func (s *CounterSink) Decrement(step int64) {
	log.Info().Msgf("sink: decrement %s: %d", s.ObjectId(), step)
	if s.node == nil {
		log.Info().Msgf("no node")
		return
	}
	methodId := core.MakeSymbolId(s.ObjectId(), "decrement")
	log.Info().Msgf("%s: %d", methodId, step)
	s.node.InvokeRemote(methodId, core.Args{step}, nil)
}

func (s *CounterSink) HandleInit(objectId string, props core.KWArgs, node *client.Node) {
	fmt.Printf("sink: on init: %s %v\n", objectId, props)
	if objectId == s.ObjectId() {
		s.node = node
		if count, ok := props["count"]; ok {
			s.count = core.AsInt(count)
		}
	}
}

func (s *CounterSink) HandlePropertyChange(propertyId string, value core.Any) {
	fmt.Printf("on property change: %s %v\n", propertyId, value)
	name := core.SymbolIdToMember(propertyId)
	switch name {
	case "count":
		s.count = core.AsInt(value)
	default:
		fmt.Printf("unknown property: %s\n", propertyId)
	}
}

func (s *CounterSink) HandleRelease() {
	fmt.Printf("on release: %s\n", s.ObjectId())
	if s.node != nil {
		s.node = nil
	}
}

func (s *CounterSink) HandleSignal(signalId string, args core.Args) {
	fmt.Printf("on signal: %s %v\n", signalId, args)
}

func main() {
	flag.Parse()
	registry := client.NewRegistry()
	ctx := context.Background()
	conn, err := ws.Dial(ctx, *addr)
	if err != nil {
		fmt.Printf("dial error: %s\n", err)
		return
	}
	defer conn.Close()
	node := client.NewNode(registry)
	node.SetOutput(conn)
	conn.SetOutput(node)
	registry.AttachClientNode(node)
	sink := NewSink(node)
	registry.AddObjectSink(sink)
	node.LinkRemoteNode(sink.ObjectId())

	if err != nil {
		fmt.Printf("dial error: %s\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sink.Increment(1)
		sink.Decrement(1)
		time.Sleep(time.Second)
	}()
	wg.Wait()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt // wait for interrupt
}
