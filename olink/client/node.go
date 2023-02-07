package client

import (
	"fmt"
	"io"

	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var nodeId = 0

func nextNodeId() string {
	nodeId++
	return fmt.Sprintf("n%d", nodeId)
}

type InvokeReplyArg struct {
	Identifier string
	Value      core.Any
}

type InvokeReplyFunc func(arg InvokeReplyArg)

type Node struct {
	id       string
	Registry *Registry
	pending  map[int64]InvokeReplyFunc
	seqId    int64
	conv     core.MessageConverter
	output   io.WriteCloser
}

func NewNode(registry *Registry) *Node {
	return &Node{
		id:       nextNodeId(),
		Registry: registry,
		pending:  make(map[int64]InvokeReplyFunc),
		seqId:    0,
		conv: core.MessageConverter{
			Format: core.FormatJson,
		},
	}
}

func (n *Node) Id() string {
	return n.id
}

func (n *Node) Close() error {
	log.Debug().Msgf("node %s: closing", n.Id())
	n.Registry.DetachClientNode(n)
	return nil
}

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
		sink := n.Registry.ObjectSink(objectId)
		if sink == nil {
			return 0, fmt.Errorf("no sink for %s", objectId)
		}
		sink.OnInit(objectId, props, n)
		return 0, nil
	case core.MsgPropertyChange:
		// get the sink and call the on property change method
		propertyId, value := msg.AsPropertyChange()
		objectId := core.SymbolIdToObjectId(propertyId)
		sink := n.Registry.ObjectSink(objectId)
		if sink == nil {
			return 0, fmt.Errorf("no sink for %s", propertyId)
		}
		sink.OnPropertyChange(propertyId, value)
	case core.MsgInvokeReply:
		// lookup the pending invoke and call the function
		requestId, methodId, value := msg.AsInvokeReply()
		log.Debug().Msgf("invoke reply: %d %s %v", requestId, methodId, value)
		fn, ok := n.pending[requestId]
		if !ok {
			return 0, fmt.Errorf("no pending invoke with id %d", requestId)
		}
		if fn == nil {
			return 0, fmt.Errorf("no function for pending invoke with id %d", requestId)
		}
		delete(n.pending, requestId)
		fn(InvokeReplyArg{methodId, value})
	case core.MsgSignal:
		// get the sink and call the on signal method
		signalId, args := msg.AsSignal()
		objectId := core.SymbolIdToObjectId(signalId)
		sink := n.Registry.ObjectSink(objectId)
		if sink == nil {
			return 0, fmt.Errorf("no sink for %s", signalId)
		}
		sink.OnSignal(signalId, args)
	case core.MsgError:
		// report the error
		msgType, id, err := msg.AsError()
		log.Info().Msgf("error: %d %d %s", msgType, id, err)
	default:
		return 0, fmt.Errorf("unknown type in client message: %#v", msg)
	}
	return len(data), nil
}

func (n *Node) InvokeRemote(methodId string, args core.Args, f InvokeReplyFunc) {
	n.seqId++
	if f != nil {
		n.pending[n.seqId] = f
	}
	n.SendMessage(core.MakeInvokeMessage(n.seqId, methodId, args))
}

func (n *Node) SetRemoteProperty(propertyId string, value core.Any) {
	n.SendMessage(core.MakeSetPropertyMessage(propertyId, value))
}

func (n *Node) LinkRemoteNode(objectId string) {
	n.Registry.LinkClientNode(objectId, n)
	n.SendMessage(core.MakeLinkMessage(objectId))
}

func (n *Node) UnlinkRemoteNode(objectId string) {
	n.Registry.UnlinkClientNode(objectId)
	n.SendMessage(core.MakeUnlinkMessage(objectId))
}
