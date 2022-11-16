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
	log.Infof("%s close\n", n.Id())
	// close(n.incoming)
	// n.incoming = nil
	return nil
}

func (n *Node) Write(data []byte) (int, error) {
	log.Debugf("node: write incoming<- %s\n", data)
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
			log.Infof("error converting message")
		}
		log.Infof("%s <- %v\n", n.Id(), msg)
		switch msg.Type() {
		case core.MsgLink:
			objectId := msg.AsLink()
			n.Registry.LinkRemoteNode(objectId, n)
			s := n.Registry.GetObjectSource(objectId)
			if s == nil {
				log.Infof("error getting object source: %s", objectId)
				break
			}
			s.Linked(objectId, n)
			// send back an init message

			props, err := s.CollectProperties()
			if err != nil {
				log.Infof("error collecting properties")
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
				log.Infof("no source for %s from %s", objectId, propertyId)
				break
			}
			s.SetProperty(propertyId, value)
			// send back property change message
			msg := core.MakePropertyChangeMessage(propertyId, value)
			n.SendMessage(msg)
		case core.MsgInvoke:
			// invoke the method on the source
			requestId, methodId, args := msg.AsInvoke()
			log.Infof("node(%s): invoke: %d %s %v\n", n.Id(), requestId, methodId, args)
			objectId := core.ToObjectId(methodId)
			s := n.Registry.GetObjectSource(objectId)
			if s == nil {
				log.Infof("no source for %s from %s", objectId, methodId)
				break
			}
			result, err := s.Invoke(methodId, args)
			if err != nil {
				log.Infof("error invoking %s: %s", methodId, err)
			}
			// send back the result
			msg := core.MakeInvokeReplyMessage(requestId, methodId, result)
			n.SendMessage(msg)
		default:
			log.Infof("unknown type in remote message: %#v type=%d", msg, msg.Type())
		}
	}
}

func (n *Node) SendMessage(msg core.Message) {
	log.Infof("%s -> %v\n", n.Id(), msg)
	if n.output == nil {
		log.Infof("error: no output for %s", n.Id())
		return
	}
	data, err := n.Converter.ToData(msg)
	if err != nil {
		log.Infof("error converting message")
		return
	}
	log.Debugf("node: write output<- %s\n", data)
	_, err = n.output.Write(data)
	if err != nil {
		log.Infof("error writing message")
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
