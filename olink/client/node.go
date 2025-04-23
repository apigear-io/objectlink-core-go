package client

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/apigear-io/objectlink-core-go/helper"
	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var nextNodeId = helper.MakeIdGenerator("n")

type InvokeReplyArg struct {
	Identifier string
	Value      core.Any
}

type InvokeReplyFunc func(arg InvokeReplyArg)

type Node struct {
	mu       sync.RWMutex
	id       string
	registry *Registry
	pending  map[int64]InvokeReplyFunc
	seqId    atomic.Int64
	conv     core.MessageConverter
	output   io.WriteCloser
}

func NewNode(registry *Registry) *Node {
	return &Node{
		id:       nextNodeId(),
		registry: registry,
		pending:  make(map[int64]InvokeReplyFunc),
		conv: core.MessageConverter{
			Format: core.FormatJson,
		},
	}
}

func (n *Node) Id() string {
	return n.id
}

func (n *Node) Registry() *Registry {
	return n.registry
}

func (n *Node) Close() error {
	log.Debug().Msgf("node %s: closing", n.Id())
	n.registry.DetachClientNode(n)
	return nil
}

// SetOutput sets the output for the node.
func (n *Node) SetOutput(out io.WriteCloser) {
	n.output = out
}

func (n *Node) SendMessage(msg core.Message) {
	log.Debug().Msgf("%s -> %v", n.Id(), msg)
	if n.output == nil {
		log.Warn().Msgf("node %s: no output", n.Id())
		return
	}
	data, err := n.conv.ToData(msg)
	if err != nil {
		log.Warn().Msgf("node %s: error converting message to data: %v", n.Id(), err)
		return
	}
	if n.output == nil {
		log.Warn().Msgf("node %s: no output", n.Id())
		return
	}
	_, err = n.output.Write(data)
	if err != nil {
		log.Warn().Msgf("node %s: error writing message: %v", n.Id(), err)
		return
	}
}

// Write handles a message from the source.
// We handle init, property change, invoke reply, signal messages.
func (n *Node) Write(data []byte) (int, error) {
	msg, err := n.conv.FromData(data)
	log.Debug().Msgf("%s <- %v", n.Id(), msg)
	if err != nil {
		return 0, err
	}
	switch msg.Type() {
	case core.MsgInit:
		// get the sink and call the on init method
		objectId, props := msg.AsInit()
		sink := n.registry.ObjectSink(objectId)
		if sink == nil {
			return 0, fmt.Errorf("no sink for %s", objectId)
		}
		sink.HandleInit(objectId, props, n)
		return 0, nil
	case core.MsgPropertyChange:
		// get the sink and call the on property change method
		propertyId, value := msg.AsPropertyChange()
		objectId := core.SymbolIdToObjectId(propertyId)
		sink := n.registry.ObjectSink(objectId)
		if sink == nil {
			return 0, fmt.Errorf("no sink for %s", propertyId)
		}
		sink.HandlePropertyChange(propertyId, value)
	case core.MsgInvokeReply:
		// lookup the pending invoke and call the function
		requestId, methodId, value := msg.AsInvokeReply()
		log.Debug().Msgf("invoke reply: %d %s %v", requestId, methodId, value)
		n.mu.RLock()
		fn, ok := n.pending[requestId]
		n.mu.RUnlock()
		if !ok {
			return 0, fmt.Errorf("no pending invoke with id %d", requestId)
		}
		if fn == nil {
			return 0, fmt.Errorf("no function for pending invoke with id %d", requestId)
		}
		n.mu.Lock()
		delete(n.pending, requestId)
		n.mu.Unlock()
		fn(InvokeReplyArg{methodId, value})
	case core.MsgSignal:
		// get the sink and call the on signal method
		signalId, args := msg.AsSignal()
		objectId := core.SymbolIdToObjectId(signalId)
		sink := n.registry.ObjectSink(objectId)
		if sink == nil {
			return 0, fmt.Errorf("no sink for %s", signalId)
		}
		sink.HandleSignal(signalId, args)
	case core.MsgError:
		// report the error
		msgType, id, err := msg.AsError()
		log.Info().Msgf("msg error: msgType=%d id-%d err=%s", msgType, id, err)
	default:
		return 0, fmt.Errorf("unknown type in client message: %#v", msg)
	}
	return len(data), nil
}

func (n *Node) InvokeRemote(methodId string, args core.Args, f InvokeReplyFunc) {
	seqId := n.seqId.Add(1)
	n.mu.Lock()
	if f != nil {
		n.pending[seqId] = f
	}
	n.mu.Unlock()
	n.SendMessage(core.MakeInvokeMessage(seqId, methodId, args))
}

func (n *Node) InvokeRemoteSync(methodId string, args core.Args) (core.Any, error) {
	ch := make(chan InvokeReplyArg, 1)
	n.InvokeRemote(methodId, args, func(arg InvokeReplyArg) {
		ch <- arg
	})
	arg := <-ch
	return arg.Value, nil
}

func (n *Node) SetRemoteProperty(propertyId string, value core.Any) {
	n.SendMessage(core.MakeSetPropertyMessage(propertyId, value))
}

func (n *Node) LinkRemoteNode(objectId string) {
	n.registry.LinkClientNode(objectId, n)
	n.SendMessage(core.MakeLinkMessage(objectId))
}

func (n *Node) UnlinkRemoteNode(objectId string) {
	n.registry.UnlinkClientNode(objectId)
	n.SendMessage(core.MakeUnlinkMessage(objectId))
}
