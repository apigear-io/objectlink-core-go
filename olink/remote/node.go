package remote

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var id = 0

func nextId() string {
	id++
	return fmt.Sprintf("n%d", id)
}

type Node struct {
	sync.Mutex
	id        string
	Registry  *Registry
	Converter core.MessageConverter
	output    io.WriteCloser
	incoming  chan []byte
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewNode(registry *Registry) *Node {
	ctx, cancel := context.WithCancel(context.Background())
	n := &Node{
		id:       nextId(),
		Registry: registry,
		Converter: core.MessageConverter{
			Format: core.FormatJson,
		},
		incoming: make(chan []byte),
		ctx:      ctx,
		cancel:   cancel,
	}
	registry.AttachRemoteNode(n)
	go n.IncomingPump()
	return n
}

func (n *Node) Id() string {
	return n.id
}

func (n *Node) SetOutput(out io.WriteCloser) {
	n.Lock()
	defer n.Unlock()
	n.output = out
}

func (n *Node) RemoveNode() {
	n.Registry.DetachRemoteNode(n)
}

func (n *Node) Close() error {
	n.cancel()
	return nil
}

func (n *Node) Write(data []byte) (int, error) {
	n.incoming <- data
	return len(data), nil
}

func (n *Node) IncomingPump() {
	for {
		select {
		case <-n.ctx.Done():
			return
		case data := <-n.incoming:
			msg, err := n.Converter.FromData(data)
			if err != nil {
				continue
			}
			switch msg.Type() {
			case core.MsgLink:
				objectId := msg.AsLink()
				n.Registry.LinkRemoteNode(objectId, n)
				s := n.Registry.GetObjectSource(objectId)
				if s == nil {
					break
				}
				s.Linked(objectId, n)
				// send back an init message

				props, err := s.CollectProperties()
				if err != nil {
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
				objectId, name := core.SymbolIdToParts(propertyId)
				s := n.Registry.GetObjectSource(objectId)
				if s == nil {
					break
				}
				s.SetProperty(name, value)
				// send back property change message
				msg := core.MakePropertyChangeMessage(propertyId, value)
				n.SendMessage(msg)
			case core.MsgInvoke:
				// invoke the method on the source
				requestId, methodId, args := msg.AsInvoke()
				objectId, name := core.SymbolIdToParts(methodId)
				s := n.Registry.GetObjectSource(objectId)
				if s == nil {
					log.Debug().Msgf("node: no source for %s", objectId)
					break
				}
				result, err := s.Invoke(name, args)
				if err != nil {
					log.Debug().Msgf("node: error invoking %s: %v", methodId, err)
					msg := core.MakeErrorMessage(core.MsgInvoke, requestId, err.Error())
					n.SendMessage(msg)
					break
				}
				log.Debug().Msgf("node: invoke result: %v", result)
				msg := core.MakeInvokeReplyMessage(requestId, methodId, result)
				n.SendMessage(msg)
			default:
				log.Info().Msgf("node: unknown message type: %v", msg.Type())
			}
		}
	}
}

func (n *Node) SendMessage(msg core.Message) {
	log.Debug().Msgf("%s -> %v", n.Id(), msg)
	n.Lock()
	out := n.output
	n.Unlock()
	if out == nil {
		log.Info().Msgf("node: no output")
		return
	}
	data, err := n.Converter.ToData(msg)
	if err != nil {
		log.Info().Msgf("node: error converting message: %v", err)
		return
	}
	log.Debug().Msgf("node: write output-> %s", data)
	_, err = out.Write(data)
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
	objectId := core.SymbolIdToObjectId(propertyId)
	msg := core.MakePropertyChangeMessage(propertyId, value)
	n.BroadcastMessage(objectId, msg)
}

func (n *Node) NotifySignal(signalId string, args core.Args) {
	objectId := core.SymbolIdToObjectId(signalId)
	msg := core.MakeSignalMessage(signalId, args)
	n.BroadcastMessage(objectId, msg)
}
