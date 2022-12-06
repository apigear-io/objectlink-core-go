package client

import (
	"fmt"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"log"
)

var (
	ErrNoWriter = fmt.Errorf("no writer")
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
	writer    core.MessageWriter
}

func NewNode(registry *Registry, writer core.MessageWriter) *Node {
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

func (n *Node) WriteMessage(msg core.Message) {
	if n.writer == nil {
		log.Printf("no writer")
		return
	}
	err := n.writer.WriteMessage(msg)
	if err != nil {
		log.Printf("error writing message: %s", err)
	}
}

// HandleMessage handles a message from the source.
// We handle init, property change, invoke reply, signal messages.
func (n *Node) HandleMessage(msg core.Message) error {
	switch msg.Type() {
	case core.MsgInit:
		// get the sink and call the on init method
		name, props := msg.AsInit()
		sink := n.Registry.GetObjectSink(name)
		if sink == nil {
			return fmt.Errorf("no sink for %s", name)
		}
		sink.OnInit(name, props, n)
		return nil
	case core.MsgPropertyChange:
		// get the sink and call the on property change method
		res, value := msg.AsPropertyChange()
		sink := n.Registry.GetObjectSink(res.ObjectId())
		if sink == nil {
			return fmt.Errorf("no sink for %s", res)
		}
		sink.OnPropertyChange(res, value)
	case core.MsgInvokeReply:
		// lookup the pending invoke and call the function
		id, res, value := msg.AsInvokeReply()
		f, ok := n.pending[id]
		if !ok {
			return fmt.Errorf("no pending invoke with id %d", id)
		}
		delete(n.pending, id)
		f(InvokeReplyArg{res, value})
	case core.MsgSignal:
		// get the sink and call the on signal method
		res, args := msg.AsSignal()
		sink := n.Registry.GetObjectSink(res.ObjectId())
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

func (c *Node) LinkRemoteNode(name string) {
	c.WriteMessage(core.CreateLinkMessage(name))
}

func (c *Node) UnlinkRemoteNode(name string) {
	c.WriteMessage(core.CreateUnlinkMessage(name))
}
