package client

import (
	"fmt"
	"log"
	"olink/pkg/core"
)

type ClientNode struct {
	invokesPending map[int]InvokeReplyFunc
	requestId      int
	core.BaseNode
}

func NewClientNode() *ClientNode {
	return &ClientNode{
		invokesPending: make(map[int]InvokeReplyFunc),
		requestId:      0,
		BaseNode:       core.BaseNode{},
	}
}

func (c *ClientNode) InvokeRemote(name string, args core.Args, f InvokeReplyFunc) {
	c.requestId++
	c.invokesPending[c.requestId] = f
	c.EmitWrite(core.NewInvokeMessage(c.requestId, name, args))
}

func (c *ClientNode) SetRemoteProperty(name string, value core.Any) {
	c.EmitWrite(core.NewSetPropertyMessage(name, value))
}

func (c *ClientNode) Registry() *ClientRegistry {
	return GetRegistry()
}

func (c *ClientNode) LinkNode(name string) {
	c.Registry().LinkClientNode(name, c)
}

func (c *ClientNode) UnlinkNode(name string) {
	c.Registry().UnlinkClientNode(name)
}

func (c *ClientNode) AddObjectSink(sink IObjectSink) {
	GetRegistry().AddObjectSink(sink)
}

func (c *ClientNode) RemoveObjectSink(sink IObjectSink) {
	GetRegistry().RemoveObjectSink(sink)
}

func (c *ClientNode) GetObjectSink(name string) IObjectSink {
	return GetRegistry().GetObjectSink(name)
}

func (c *ClientNode) LinkRemote(name string) {
	c.LinkNode(name)
	c.EmitWrite(core.NewLinkMessage(name))
}

func (c *ClientNode) UnlinkRemote(name string) {
	c.UnlinkNode(name)
	c.EmitWrite(core.NewUnlinkMessage(name))
}

func (c *ClientNode) HandleInit(name string, props core.Props) error {
	sink := c.GetObjectSink(name)
	if sink == nil {
		sink.OnInit(name, props, c)
	} else {
		return fmt.Errorf("no sink for %s", name)
	}
	return nil
}

func (c *ClientNode) HandlePropertyChange(name string, value core.Any) error {
	sink := c.GetObjectSink(name)
	if sink == nil {
		sink.OnPropertyChange(name, value)
	} else {
		return fmt.Errorf("no sink for %s", name)
	}
	return nil

}

func (c *ClientNode) HandleInvokeReply(id int, name string, value core.Any) error {
	f, ok := c.invokesPending[id]
	if ok {
		delete(c.invokesPending, id)
		f(InvokeReplyArg{name, value})
	} else {
		return fmt.Errorf("no pending invoke with id %d", id)
	}
	return nil
}

func (c *ClientNode) HandleSignal(name string, args core.Args) error {
	sink := c.GetObjectSink(name)
	if sink == nil {
		sink.OnSignal(name, args)
	} else {
		return fmt.Errorf("no sink for %s", name)
	}
	return nil
}

func (c *ClientNode) HandleError(msgType core.MsgType, id int, error string) error {
	log.Printf("client node error: %s", error)
	return nil
}
