package client

import (
	"fmt"

	"github.com/apigear-io/objectlink-core-go/log"
)

type SinkFactory func(objectId string) IObjectSink

type SinkToClientEntry struct {
	sink IObjectSink
	node *Node
}

var registryId = 0

func nextRegistryId() string {
	registryId++
	return fmt.Sprintf("r%d", registryId)
}

// Registry is a registry of object sinks.
// It is used to keep track of object sinks and their associated client nodes.
// It is optimized for the retrieval of object sinks by object id.
// A sink is always associated with zero or one client node.
// A node can be linked to zero or many sinks.
type Registry struct {
	id      string
	entries map[string]*SinkToClientEntry
	factory SinkFactory
}

func NewRegistry() *Registry {
	return &Registry{
		id:      nextRegistryId(),
		entries: make(map[string]*SinkToClientEntry),
	}
}

func (r *Registry) Id() string {
	return r.id
}

// SetSinkFactory sets the sink factory.
func (r *Registry) SetSinkFactory(factory SinkFactory) {
	r.factory = factory
}

// attach client node to registry
func (r *Registry) AttachClientNode(node *Node) {
}

// detach client node from registry
func (r *Registry) DetachClientNode(node *Node) {
	for _, e := range r.entries {
		if e.node == node {
			log.Debug().Msgf("unlink client node %s from object %s", node.Id(), e.sink.ObjectId())
			e.node = nil
		}
	}
}

func (r *Registry) LinkClientNode(objectId string, node *Node) {
	if entry := r.Entry(objectId); entry != nil {
		log.Debug().Msgf("link client node %s to object %s", node.Id(), objectId)
		entry.node = node
	} else {
		log.Warn().Msgf("object %s not found", objectId)
	}
}

func (r *Registry) UnlinkClientNode(objectId string) {
	log.Debug().Msgf("unlink client node from object %s", objectId)
	r.Entry(objectId).node = nil
}

func (r *Registry) AddObjectSink(sink IObjectSink) {
	log.Info().Msgf("add object sink %s", sink.ObjectId())
	r.Entry(sink.ObjectId()).sink = sink
}

// remove object sink from registry
func (r *Registry) RemoveObjectSink(objectId string) {
	log.Info().Msgf("remove object sink %s", objectId)
	sink := r.Entry(objectId).sink

	if sink != nil {
		sink.OnRelease()
	} else {
		log.Warn().Msgf("object sink %s not found", objectId)
	}
	r.RemoveEntry(objectId)
}

// get object sink by name
func (r *Registry) ObjectSink(objectId string) IObjectSink {
	s := r.Entry(objectId).sink
	if s == nil && r.factory != nil {
		s = r.factory(objectId)
		r.Entry(objectId).sink = s
	}
	return s
}

func (r *Registry) Node(objectId string) *Node {
	return r.Entry(objectId).node
}

func (r *Registry) Entry(objectId string) *SinkToClientEntry {
	if r.entries[objectId] == nil {
		r.entries[objectId] = &SinkToClientEntry{
			node: nil,
			sink: nil,
		}
	}
	return r.entries[objectId]
}

func (r *Registry) RemoveEntry(objectId string) {
	delete(r.entries, objectId)
}

func (r *Registry) ObjectIds() []string {
	ids := make([]string, 0, len(r.entries))
	for id := range r.entries {
		ids = append(ids, id)
	}
	return ids
}
