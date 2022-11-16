package remote

import (
	"fmt"
	"io"
	"olink/log"
	"olink/pkg/core"
)

var id = 0

func nextId() string {
	id++
	return fmt.Sprintf("n%d", id)
}

type Node struct {
	id        string
	Registry  *Registry
	Converter core.MessageConverter
	output    io.WriteCloser
	incoming  chan []byte
}

func NewNode(registry *Registry) *Node {
	n := &Node{
		id:       nextId(),
		Registry: registry,
		Converter: core.MessageConverter{
			Format: core.FormatJson,
		},
		incoming: make(chan []byte),
	}
	registry.AttachRemoteNode(n)
	go n.IncomingPump()
	return n
}

func (n *Node) Id() string {
	return n.id
}

func (n *Node) SetOutput(out io.WriteCloser) {
	n.output = out
}

func (n *Node) RemoveNode() {
	n.Registry.DetachRemoteNode(n)
}

func (n *Node) Close() error {
	log.Info().Msgf("node: close %s", n.Id())
	// close(n.incoming)
	// n.incoming = nil
	return nil
}

func (n *Node) Write(data []byte) (int, error) {
	log.Debug().Msgf("node: write input-> %s\n", data)
	// if n.incoming == nil {
	// 	return 0, fmt.Errorf("node: write: node is closed")
	// }
	n.incoming <- data
	return len(data), nil
}

func (n *Node) IncomingPump() {
	for data := range n.incoming {
		msg, err := n.Converter.FromData(data)
		if err != nil {
			log.Info().Msgf("node: error parsing message: %v", err)
		}
		log.Info().Msgf("%s <- %v", n.Id(), msg)
		switch msg.Type() {
		case core.MsgLink:
			objectId := msg.AsLink()
			n.Registry.LinkRemoteNode(objectId, n)
			s := n.Registry.GetObjectSource(objectId)
			if s == nil {
				log.Info().Msgf("node: no source for %s", objectId)
				break
			}
			s.Linked(objectId, n)
			// send back an init message

			props, err := s.CollectProperties()
			if err != nil {
				log.Info().Msgf("node: error collecting properties: %v", err)
				break
			}
			msg := core.MakeInitMessage(objectId, props)
			n.SendMessage(msg)
		case core.MsgUnlink:
			// unlink the sink from the source
			objectId := msg.AsUnlink()
			n.Registry.UnlinkRemoteNode(objectId, n)
		case core.MsgSetProperty:
			// set the property on the source
			propertyId, value := msg.AsSetProperty()
			objectId := core.ToObjectId(propertyId)
			s := n.Registry.GetObjectSource(objectId)
			if s == nil {
				log.Info().Msgf("node: no source for %s", objectId)
				break
			}
			s.SetProperty(propertyId, value)
			// send back property change message
			msg := core.MakePropertyChangeMessage(propertyId, value)
			n.SendMessage(msg)
		case core.MsgInvoke:
			// invoke the method on the source
			requestId, methodId, args := msg.AsInvoke()
			log.Info().Msgf("node: invoke %d %s %v", requestId, methodId, args)
			objectId := core.ToObjectId(methodId)
			s := n.Registry.GetObjectSource(objectId)
			if s == nil {
				log.Info().Msgf("node: no source for %s", objectId)
				break
			}
			result, err := s.Invoke(methodId, args)
			if err != nil {
				log.Info().Msgf("node: error invoking method: %v", err)
			}
			// send back the result
			msg := core.MakeInvokeReplyMessage(requestId, methodId, result)
			n.SendMessage(msg)
		default:
			log.Info().Msgf("node: unknown message type: %v", msg.Type())
		}
	}
}

func (n *Node) SendMessage(msg core.Message) {
	log.Info().Msgf("%s -> %v", n.Id(), msg)
	if n.output == nil {
		log.Info().Msgf("node: no output")
		return
	}
	data, err := n.Converter.ToData(msg)
	if err != nil {
		log.Info().Msgf("node: error converting message: %v", err)
		return
	}
	log.Debug().Msgf("node: write output-> %s\n", data)
	_, err = n.output.Write(data)
	if err != nil {
		log.Info().Msgf("node: error writing message: %v", err)
	}
}

func (n *Node) BroadcastMessage(objectId string, msg core.Message) {
	for _, node := range n.Registry.GetRemoteNodes(objectId) {
		node.SendMessage(msg)
	}
}

func (n *Node) NotifyPropertyChange(propertyId string, value core.Any) {
	objectId := core.ToObjectId(propertyId)
	msg := core.MakePropertyChangeMessage(propertyId, value)
	n.BroadcastMessage(objectId, msg)
}

func (n *Node) NotifySignal(signalId string, args core.Args) {
	objectId := core.ToObjectId(signalId)
	msg := core.MakeSignalMessage(signalId, args)
	n.BroadcastMessage(objectId, msg)
}
