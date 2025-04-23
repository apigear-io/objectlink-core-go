package client

import (
	"fmt"
	"sync"

	"github.com/apigear-io/objectlink-core-go/helper"
	"github.com/apigear-io/objectlink-core-go/log"
)

type SinkFactory func(objectId string) IObjectSink

var nextRegistryId = helper.MakeIdGenerator("r")

func clearRegistryId() {
	nextRegistryId = helper.MakeIdGenerator("r")
}

// Registry is a registry of object sinks.
// It is used to keep track of object sinks and their associated client nodes.
// It is optimized for the retrieval of object sinks by object id.
// A sink is always associated with zero or one client node.
// A node can be linked to zero or many sinks.
type Registry struct {
	sync.RWMutex
	id      string
	entries *clientEntries
}

func NewRegistry() *Registry {
	log.Debug().Msg("create new registry")
	return &Registry{
		id:      nextRegistryId(),
		entries: newClientEntries(),
	}
}

func (r *Registry) Id() string {
	r.RLock()
	defer r.RUnlock()
	return r.id
}

// SetSinkFactory sets the sink factory.
func (r *Registry) SetSinkFactory(factory SinkFactory) {
	log.Debug().Msg("set sink factory")
	r.entries.setFactory(factory)
}

// attach client node to registry
func (r *Registry) AttachClientNode(node *Node) {
}

// detach client node from registry
func (r *Registry) DetachClientNode(node *Node) {
	if node == nil {
		return
	}
	log.Debug().Msgf("detach client node %s", node.Id())
	r.entries.purgeNode(node)
}

func (r *Registry) LinkClientNode(objectId string, node *Node) {
	log.Debug().Msgf("link client node to object %s", objectId)
	r.entries.setNode(objectId, node)
}

func (r *Registry) UnlinkClientNode(objectId string) {
	log.Debug().Msgf("unlink client node from object %s", objectId)
	r.entries.clearNode(objectId)
}

func (r *Registry) GetClientNode(objectId string) *Node {
	return r.entries.getNode(objectId)
}

func (r *Registry) AddObjectSink(sink IObjectSink) error {
	if sink == nil {
		return fmt.Errorf("object sink is nil")
	}
	return r.entries.setSink(sink)
}

func (r *Registry) IsRegistered(objectId string) bool {
	return r.entries.hasEntry(objectId)
}

// remove object sink from registry
func (r *Registry) RemoveObjectSink(objectId string) {
	log.Info().Msgf("remove object sink %s", objectId)
	s := r.entries.getSink(objectId)

	if s != nil {
		s.HandleRelease()
	} else {
		log.Warn().Msgf("object sink %s not found", objectId)
	}
	r.entries.removeEntry(objectId)
}

// get object sink by name
func (r *Registry) ObjectSink(objectId string) IObjectSink {
	return r.entries.getSink(objectId)
}

func (r *Registry) ObjectIds() []string {
	return r.entries.getEntryIds()
}
