package remote

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/apigear-io/objectlink-core-go/helper"
	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var nextNodeId = helper.MakeIdGenerator("n")

type Node struct {
	sync.RWMutex
	id       string
	registry *Registry
	conv     core.MessageConverter
	output   io.WriteCloser
	incoming chan []byte
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewNode(registry *Registry) *Node {
	nodeId := nextNodeId()
	log.Debug().Msgf("node %s: creating", nodeId)
	ctx, cancel := context.WithCancel(context.Background())
	n := &Node{
		id:       nodeId,
		registry: registry,
		conv: core.MessageConverter{
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
	n.RLock()
	defer n.RUnlock()
	return n.id
}

func (n *Node) Registry() *Registry {
	n.RLock()
	defer n.RUnlock()
	return n.registry
}

func (n *Node) SetOutput(out io.WriteCloser) {
	n.Lock()
	n.output = out
	n.Unlock()
}

func (n *Node) RemoveNode() {
	n.RLock()
	registry := n.registry
	n.RUnlock()
	registry.DetachRemoteNode(n)
}

func (n *Node) Close() error {
	log.Debug().Msgf("node %s: closing", n.id)
	n.RLock()
	cancel := n.cancel
	cancel()
	n.RUnlock()
	return nil
}

func (n *Node) Write(data []byte) (int, error) {
	n.RLock()
	incoming := n.incoming
	n.RUnlock()
	incoming <- data
	return len(data), nil
}

func (n *Node) IncomingPump() {
	for {
		select {
		case <-n.ctx.Done():
			return
		case data := <-n.incoming:
			msg, err := n.conv.FromData(data)
			if err != nil {
				continue
			}
			switch msg.Type() {
			case core.MsgLink:
				objectId := msg.AsLink()
				n.registry.LinkRemoteNode(objectId, n)
				s := n.registry.GetObjectSource(objectId)
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
				n.registry.UnlinkRemoteNode(objectId, n)
			case core.MsgSetProperty:
				// set the property on the source
				propertyId, value := msg.AsSetProperty()
				objectId, name := core.SymbolIdToParts(propertyId)
				s := n.registry.GetObjectSource(objectId)
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
				s := n.registry.GetObjectSource(objectId)
				if s == nil {
					log.Warn().Msgf("node: no source for %s", objectId)
					break
				}
				result, err := s.Invoke(name, args)
				if err != nil {
					log.Warn().Msgf("node: error invoking %s: %v", methodId, err)
					msg := core.MakeErrorMessage(core.MsgInvoke, requestId, err.Error())
					n.SendMessage(msg)
					break
				}
				log.Debug().Msgf("node: invoke result: %v", result)
				msg := core.MakeInvokeReplyMessage(requestId, methodId, result)
				n.SendMessage(msg)
			case core.MsgSignal:
				// send the signal to all nodes
				signalId, args := msg.AsSignal()
				if n.registry != nil {
					objectId, name := core.SymbolIdToParts(signalId)
					n.registry.NotifySignal(objectId, name, args)
				} else {
					n.SendSignal(signalId, args)
				}
			default:
				log.Info().Msgf("node: unknown message type: %v", msg.Type())
			}
		}
	}
}

func (n *Node) SendMessage(msg core.Message) {
	log.Debug().Msgf("-> %s send %v", n.id, msg)
	n.RLock()
	output := n.output
	conv := n.conv
	n.RUnlock()
	err := doSendMessage(output, conv, msg)
	if err != nil {
		log.Error().Msgf("node: error sending message: %v", err)
	}
}

func doSendMessage(o io.WriteCloser, c core.MessageConverter, msg core.Message) error {
	if o == nil {
		return fmt.Errorf("no output")
	}
	if msg == nil {
		return fmt.Errorf("no message")
	}

	data, err := c.ToData(msg)
	if err != nil {
		return fmt.Errorf("error converting message: %v", err)
	}
	_, err = o.Write(data)
	if err != nil {
		return fmt.Errorf("error writing message: %v", err)
	}
	return nil
}

func (n *Node) NotifyPropertyChange(propertyId string, value core.Any) {
	log.Debug().Msgf("node %s: notify property change: %s", n.id, propertyId)
	if n.registry != nil {
		objectId, name := core.SymbolIdToParts(propertyId)
		n.registry.NotifyPropertyChange(objectId, core.KWArgs{name: value})
	} else {
		n.SendPropertyChange(propertyId, value)
	}
}

func (n *Node) SendPropertyChange(propertyId string, value core.Any) {
	log.Debug().Msgf("node %s: send property change: %s", n.id, propertyId)
	msg := core.MakePropertyChangeMessage(propertyId, value)
	n.SendMessage(msg)
}

func (n *Node) NotifySignal(signalId string, args core.Args) {
	log.Debug().Msgf("node %s: notify signal: %s", n.id, signalId)
	objectId, name := core.SymbolIdToParts(signalId)
	if n.registry != nil {
		n.registry.NotifySignal(objectId, name, args)
	} else {
		n.SendSignal(signalId, args)
	}
}

func (n *Node) SendSignal(signalId string, args core.Args) {
	log.Debug().Msgf("node %s: send signal: %s", n.id, signalId)
	msg := core.MakeSignalMessage(signalId, args)
	n.SendMessage(msg)
}
