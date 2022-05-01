package client

import (
	"fmt"
	"log"
	"olink/pkg/core"
)

type InvokeReplyArg struct {
	resource core.Resource
	value    core.Any
}

type InvokeReplyFunc func(args InvokeReplyArg)

type Node struct {
	Registry  *Registry
	pending   map[int]InvokeReplyFunc
	seqId     int
	converter core.MessageConverter
	writer    core.DataWriter
}

func NewNode(registry *Registry, writer core.DataWriter) *Node {
	return &Node{
		Registry: registry,
		pending:  make(map[int]InvokeReplyFunc),
		seqId:    0,
		converter: core.MessageConverter{
			Format: core.FormatJson,
		},
		writer: writer,
	}
}

func (node *Node) WriteMessage(msg core.Message) {
	data, err := node.converter.ToData(msg)
	if err != nil {
		fmt.Printf("error converting message")
		return
	}
	err = node.writer.WriteData(data)
	if err != nil {
		fmt.Printf("error writing message")
		return
	}
}

// HandleMessage handles a message from the source.
// We handle init, property change, invoke reply, signal messages.
func (c *Node) HandleMessage(data []byte) error {
	msg, err := c.converter.FromData(data)
	if err != nil {
		return err
	}
	switch msg.Type() {
	case core.MsgInit:
		// get the sink and call the on init method
		name, props := msg.AsInit()
		sink := c.GetObjectSink(name)
		if sink == nil {
			return fmt.Errorf("no sink for %s", name)
		}
		sink.OnInit(name, props, c)
		return nil
	case core.MsgPropertyChange:
		// get the sink and call the on property change method
		res, value := msg.AsPropertyChange()
		sink := c.GetObjectSink(res.ObjectId())
		if sink == nil {
			return fmt.Errorf("no sink for %s", res)
		}
		sink.OnPropertyChange(res, value)
	case core.MsgInvokeReply:
		// lookup the pending invoke and call the function
		id, res, value := msg.AsInvokeReply()
		f, ok := c.pending[id]
		if !ok {
			return fmt.Errorf("no pending invoke with id %d", id)
		}
		delete(c.pending, id)
		f(InvokeReplyArg{res, value})
	case core.MsgSignal:
		// get the sink and call the on signal method
		res, args := msg.AsSignal()
		sink := c.GetObjectSink(res.ObjectId())
		if sink == nil {
			return fmt.Errorf("no sink for %s", res)
		}
		sink.OnSignal(res, args)
	case core.MsgError:
		// report the error
		msgType, id, error := msg.AsError()
		log.Printf("client node error: %d, %d, %s", msgType, id, error)
	default:
		// unknown message type
		return fmt.Errorf("unknown message type %d", msg.Type())
	}
	return nil
}

func (c *Node) InvokeRemote(res core.Resource, args core.Args, f InvokeReplyFunc) {
	c.seqId++
	c.pending[c.seqId] = f
	c.WriteMessage(core.CreateInvokeMessage(c.seqId, res, args))
}

func (c *Node) SetRemoteProperty(res core.Resource, value core.Any) {
	c.WriteMessage(core.CreateSetPropertyMessage(res, value))
}

func (c *Node) LinkNode(objectId string) {
	c.Registry.LinkClientNode(objectId, c)
}

func (c *Node) UnlinkNode(name string) {
	c.Registry.UnlinkClientNode(name)
}

func (c *Node) AddObjectSink(sink IObjectSink) {
	c.Registry.AddObjectSink(sink)
}

func (c *Node) RemoveObjectSink(sink IObjectSink) {
	c.Registry.RemoveObjectSink(sink)
}

func (c *Node) GetObjectSink(name string) IObjectSink {
	return c.Registry.GetObjectSink(name)
}

func (c *Node) LinkRemoteNode(name string) {
	c.LinkNode(name)
	c.WriteMessage(core.CreateLinkMessage(name))
}

func (c *Node) UnlinkRemoteNode(name string) {
	c.UnlinkNode(name)
	c.WriteMessage(core.CreateUnlinkMessage(name))
}
