package main

import (
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/ws"
	"log"
)

type Counter struct {
	node *client.Node
}

func NewCounter(node *client.Node) *Counter {
	return &Counter{
		node: node,
	}
}

func (c *Counter) Increment() {
	c.node.InvokeRemote(core.CreateResource("demo.Counter", "increment"), core.Args{}, nil)
}

func (c *Counter) Decrement() {
	c.node.InvokeRemote(core.CreateResource("demo.Counter", "decrement"), core.Args{}, nil)
}

func main() {
	c := ws.NewClient("ws://localhost:8080")
	err := c.Connect()
	if err != nil {
		log.Fatal(err)
	}
	r := client.NewRegistry()
	n := client.NewNode(r, c)
	counter := NewCounter(n)
	counter.Increment()
	counter.Increment()
}
