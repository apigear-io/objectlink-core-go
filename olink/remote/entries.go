package remote

import (
	"fmt"
	"sync"

	"github.com/apigear-io/objectlink-core-go/log"
)

type SourceFactory func(objectId string) IObjectSource

// remoteEntries is a map of object id to sourceToNodeEntry
type remoteEntries struct {
	sync.RWMutex
	entries map[string]*sourceToNodeEntry
	factory SourceFactory
}

// newRemoteEntries creates a new remoteEntries
func newRemoteEntries() *remoteEntries {
	return &remoteEntries{
		entries: make(map[string]*sourceToNodeEntry),
	}
}

// setFactory sets the source factory
func (r *remoteEntries) setFactory(factory SourceFactory) {
	r.Lock()
	defer r.Unlock()
	r.factory = factory
}

// hasEntry returns true if the entry exists
func (r *remoteEntries) hasEntry(objectId string) bool {
	r.RLock()
	defer r.RUnlock()
	_, ok := r.entries[objectId]
	return ok
}

// removeEntry removes the entry
func (r *remoteEntries) removeEntry(objectId string) {
	r.Lock()
	defer r.Unlock()
	delete(r.entries, objectId)
}

// addSource adds the source to the entry
// if the source already exists, it returns an error
// the entry is identified by the source's object id
func (r *remoteEntries) addSource(source IObjectSource) error {
	if source == nil {
		return fmt.Errorf("source is nil")
	}
	log.Info().Msgf("registry: add %s", source.ObjectId())
	e := r.getEntry(source.ObjectId())
	if e.hasSource() {
		return fmt.Errorf("source %s is already registered", source.ObjectId())
	}
	e.setSource(source)
	return nil
}

// getSource returns the source
// if the source does not exist, it is created using a factory
func (r *remoteEntries) getSource(objectId string) IObjectSource {
	e := r.getEntry(objectId)
	r.Lock()
	factory := r.factory
	r.Unlock()
	if !e.hasSource() && factory != nil {
		e.setSource(r.factory(objectId))
	}
	return e.getSource()
}

// addNode adds the node to the entry
func (r *remoteEntries) addNode(objectId string, node *Node) error {
	e := r.getEntry(objectId)
	return e.addNode(node)
}

// removeNode removes the node from the entry
func (r *remoteEntries) removeNode(objectId string, node *Node) {
	e := r.getEntry(objectId)
	e.removeNode(node)
}

// purgeNode removes the node from all entries
func (r *remoteEntries) purgeNode(node *Node) {
	r.RLock()
	defer r.RUnlock()
	for _, e := range r.entries {
		e.removeNode(node)
	}
}

// getNodes returns the list of nodes
func (r *remoteEntries) getNodes(objectId string) []*Node {
	e := r.getEntry(objectId)
	return e.getNodes()
}

// getEntry returns the entry
// if the entry does not exist, it is created
func (r *remoteEntries) getEntry(objectId string) *sourceToNodeEntry {
	r.Lock()
	defer r.Unlock()
	e, ok := r.entries[objectId]
	if !ok {
		e = &sourceToNodeEntry{}
		r.entries[objectId] = e
	}
	return e
}
