package client

import (
	"sync"

	"github.com/apigear-io/objectlink-core-go/log"
)

type SinkToClientEntry struct {
	sync.RWMutex
	sink IObjectSink
	node *Node
}

// setNode sets the node.
func (e *SinkToClientEntry) setNode(node *Node) {
	e.Lock()
	defer e.Unlock()
	e.node = node
}

// getNode returns the node.
func (e *SinkToClientEntry) getNode() *Node {
	e.RLock()
	defer e.RUnlock()
	return e.node
}

// clearNode clears the node.
func (e *SinkToClientEntry) clearNode() {
	e.Lock()
	defer e.Unlock()
	e.node = nil
}

// setSink sets the sink.
func (e *SinkToClientEntry) setSink(sink IObjectSink) {
	log.Debug().Msgf("setSink: %s", sink.ObjectId())
	e.Lock()
	defer e.Unlock()
	e.sink = sink
}

// getSink returns the sink.
func (e *SinkToClientEntry) getSink() IObjectSink {
	e.RLock()
	defer e.RUnlock()
	return e.sink
}

// hasSink returns true if the sink is set.
func (e *SinkToClientEntry) hasSink() bool {
	e.RLock()
	defer e.RUnlock()
	return e.sink != nil
}

// clearSink clears the sink.
func (e *SinkToClientEntry) clearSink() {
	e.Lock()
	defer e.Unlock()
	e.sink = nil
}
