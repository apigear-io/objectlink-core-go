package remote

import (
	"fmt"

	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/apigear-io/objectlink-core-go/olink/core"
)

type SourceFactory func(objectId string) IObjectSource

var registryId = 0

func nextRegistryId() string {
	registryId++
	return fmt.Sprintf("r%d", registryId)
}

type SourceToNodeEntry struct {
	source IObjectSource
	nodes  []*Node
}

// Registry is the registry of remote objects.
// It is optimized for the retrieval of object sources
// A object source is registered in the registry and can be retrieved by the object id.
// The source can have one or more remote nodes linked to it.
type Registry struct {
	id            string
	entries       map[string]*SourceToNodeEntry
	sourceFactory SourceFactory
}

func NewRegistry() *Registry {
	return &Registry{
		id:      nextRegistryId(),
		entries: make(map[string]*SourceToNodeEntry),
	}
}

func (r *Registry) Id() string {
	return r.id
}

// SetSourceFactory sets the source factory.
func (r *Registry) SetSourceFactory(factory SourceFactory) {
	r.sourceFactory = factory
}

// AddObjectSource adds the object source to the registry.
func (r *Registry) AddObjectSource(source IObjectSource) {
	r.entry(source.ObjectId()).source = source
}

// RemoveObjectSource removes the object source from the registry.
func (r *Registry) RemoveObjectSource(source IObjectSource) {
	r.removeEntry(source.ObjectId())
}

// GetObjectSource returns the object source by name.
func (r *Registry) GetObjectSource(objectId string) IObjectSource {
	s := r.entry(objectId).source
	if s == nil && r.sourceFactory != nil {
		s = r.sourceFactory(objectId)
		if s != nil {
			r.AddObjectSource(s)
		}
	}
	return s
}

// Checks if the object is registered.
func (r *Registry) IsRegistered(objectId string) bool {
	_, ok := r.entries[objectId]
	return ok
}

// GetRemoteNode returns the node that is linked to the object.
func (r *Registry) GetRemoteNodes(objectId string) []*Node {
	return r.entry(objectId).nodes
}

// AttachRemoteNode attaches the node to the registry.
func (r *Registry) AttachRemoteNode(node *Node) {
}

// DetachRemoteNode removes the link between the object and the node.
func (r *Registry) DetachRemoteNode(node *Node) {
	for _, v := range r.entries {
		if v.nodes != nil {
			for i, n := range v.nodes {
				if n == node {
					v.nodes = append(v.nodes[:i], v.nodes[i+1:]...)
				}
			}
		}
	}
}

// LinkRemoteNode adds a link between the object and the node.
func (r *Registry) LinkRemoteNode(objectId string, node *Node) {
	log.Info().Msgf("registry: link %s -> %s", objectId, node.Id())
	r.entry(objectId).nodes = append(r.entry(objectId).nodes, node)
}

// UnlinkRemoteNode removes the link between the object and the node.
func (r *Registry) UnlinkRemoteNode(objectId string, node *Node) {
	for i, n := range r.entry(objectId).nodes {
		if n == node {
			r.entry(objectId).nodes = append(r.entry(objectId).nodes[:i], r.entry(objectId).nodes[i+1:]...)
		}
	}
}

func (r *Registry) entry(objectId string) *SourceToNodeEntry {
	e, ok := r.entries[objectId]
	if !ok {
		e = &SourceToNodeEntry{
			source: nil,
			nodes:  make([]*Node, 0),
		}
		r.entries[objectId] = e
	}
	return e
}

func (r *Registry) removeEntry(objectId string) {
	delete(r.entries, objectId)
}

func (e *Registry) NotifyPropertyChange(objectId string, value core.Any) {
	for _, n := range e.entry(objectId).nodes {
		n.NotifyPropertyChange(objectId, value)
	}
}
