package client

import (
	"fmt"
	"sync"

	"github.com/apigear-io/objectlink-core-go/log"
)

type clientEntries struct {
	sync.RWMutex
	entries map[string]*SinkToClientEntry
	factory SinkFactory
}

func newClientEntries() *clientEntries {
	return &clientEntries{
		entries: make(map[string]*SinkToClientEntry),
	}
}

// setFactory sets the sink factory.
func (e *clientEntries) setFactory(factory SinkFactory) {
	e.Lock()
	defer e.Unlock()
	e.factory = factory
}

// hasEntry returns true if the entry exists.
func (e *clientEntries) hasEntry(objectId string) bool {
	e.RLock()
	defer e.RUnlock()
	_, ok := e.entries[objectId]
	return ok
}

// removeEntry removes the entry.
func (e *clientEntries) removeEntry(objectId string) {
	e.Lock()
	defer e.Unlock()
	delete(e.entries, objectId)
}

// getEntry returns the entry.
// if the entry does not exist, it is created using a factory.
func (e *clientEntries) getEntry(objectId string) *SinkToClientEntry {
	e.Lock()
	defer e.Unlock()
	entry, ok := e.entries[objectId]
	if !ok {
		entry = &SinkToClientEntry{}
		e.entries[objectId] = entry
	}
	return entry
}

// setSink sets the sink.
func (e *clientEntries) setSink(sink IObjectSink) error {
	entry := e.getEntry(sink.ObjectId())
	if entry.hasSink() {
		return fmt.Errorf("sink already exists for %s", sink.ObjectId())
	}
	entry.setSink(sink)
	return nil
}

// getSink returns the sink.
// if the sink does not exist, it is created using a factory.
func (e *clientEntries) getSink(objectId string) IObjectSink {
	entry := e.getEntry(objectId)
	if entry.hasSink() {
		return entry.getSink()
	}
	e.Lock()
	factory := e.factory
	e.Unlock()
	if factory != nil {
		log.Debug().Msgf("client factory: create sink for %s", objectId)
		sink := factory(objectId)
		entry.setSink(sink)
		return sink
	}
	return nil
}

// purgeNode removes all entries associated with the node.
func (e *clientEntries) purgeNode(node *Node) {
	e.RLock()
	defer e.RUnlock()
	for _, entry := range e.entries {
		if entry.node == node {
			entry.clearNode()
		}
	}
}

// setNode sets the node for the object id.
func (e *clientEntries) setNode(objectId string, node *Node) {
	entry := e.getEntry(objectId)
	entry.setNode(node)
}

// getNode returns the node for the object id.
func (e *clientEntries) getNode(objectId string) *Node {
	entry := e.getEntry(objectId)
	return entry.getNode()
}

// clearNode clears the node for the object id.
func (e *clientEntries) clearNode(objectId string) {
	entry := e.getEntry(objectId)
	entry.clearNode()
}

// getEntryIds returns the object ids.
func (e *clientEntries) getEntryIds() []string {
	e.RLock()
	defer e.RUnlock()
	ids := make([]string, 0, len(e.entries))
	for id := range e.entries {
		ids = append(ids, id)
	}
	return ids
}
