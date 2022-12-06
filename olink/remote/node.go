package remote

import (
	"fmt"
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

type Node struct {
	registry *Registry
	writer   core.MessageWriter
}

func NewNode(r *Registry, w core.MessageWriter) *Node {
	return &Node{
		registry: r,
		writer:   w,
	}
}

func (n *Node) RemoveNode() {
	n.registry.DetachRemoteNode(n)
}

func (n *Node) HandleMessage(msg core.Message) error {
	switch msg.Type() {
	case core.MsgLink:
		objectId := msg.AsLink()
		n.registry.LinkRemoteNode(objectId, n)
		s := n.registry.GetObjectSource(objectId)
		if s == nil {
			return fmt.Errorf("no source for %s", objectId)
		}
		s.Linked(objectId, n)
		// send back an init message

		props := s.CollectProperties()
		msg := core.CreateInitMessage(objectId, props)
		n.WriteMessage(msg)

	case core.MsgUnlink:
		// unlink the sink from the source
		objectId := msg.AsUnlink()
		n.registry.UnlinkRemoteNode(objectId, n)
	case core.MsgSetProperty:
		// set the property on the source
		res, value := msg.AsSetProperty()
		s := n.registry.GetObjectSource(res.ObjectId())
		if s == nil {
			return fmt.Errorf("no source for %s", res)
		}
		s.SetProperty(res, value)
		// send back property change message
		msg := core.CreatePropertyChangeMessage(res, value)
		n.WriteMessage(msg)
	default:
		return fmt.Errorf("unhandled message type %d", msg.Type())
	}
	return nil
}

func (n *Node) WriteMessage(msg core.Message) {
	if n.writer == nil {
		return
	}
	err := n.writer.WriteMessage(msg)
	if err != nil {
		fmt.Printf("error writing message")
		return
	}
}
