package remote

import (
	"fmt"
	"olink/pkg/core"
)

type Node struct {
	Registry  *Registry
	Converter *core.MessageConverter
	Writer    core.DataWriter
}

func NewNode(registry *Registry) *Node {
	return &Node{
		Registry: registry,
		Converter: &core.MessageConverter{
			Format: core.FormatJson,
		},
	}
}

func (n *Node) RemoveNode() {
	n.Registry.DetachRemoteNode(n)
}

func (n *Node) HandleMessage(data []byte) error {
	msg, err := n.Converter.FromData(data)
	if err != nil {
		return err
	}
	switch msg.Type() {
	case core.MsgLink:
		objectId := msg.AsLink()
		n.Registry.LinkRemoteNode(objectId, n)
		s := n.Registry.GetObjectSource(objectId)
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
		n.Registry.UnlinkRemoteNode(objectId, n)
	case core.MsgSetProperty:
		// set the property on the source
		res, value := msg.AsSetProperty()
		s := n.Registry.GetObjectSource(res.ObjectId())
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
	data, err := n.Converter.ToData(msg)
	if err != nil {
		fmt.Printf("error converting message")
		return
	}
	err = n.Writer.WriteData(data)
	if err != nil {
		fmt.Printf("error writing message")
		return
	}
}
