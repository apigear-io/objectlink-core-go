package remote

import (
	"github.com/apigear-io/objectlink-core-go/helper"
	"github.com/apigear-io/objectlink-core-go/log"

	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var nextRegistryId = helper.MakeIdGenerator("r")

// Registry is the registry of remote objects.
// It is optimized for the retrieval of object sources
// A object source is registered in the registry and can be retrieved by the object id.
// The source can have one or more remote nodes linked to it.
type Registry struct {
	id      string
	entries *remoteEntries
}

// NewRegistry creates a new registry.
func NewRegistry() *Registry {
	r := &Registry{
		id:      nextRegistryId(),
		entries: newRemoteEntries(),
	}
	return r
}

// Id returns the registry id.
func (r *Registry) Id() string {
	return r.id
}

// SetSourceFactory sets the source factory.
func (r *Registry) SetSourceFactory(factory SourceFactory) {
	r.entries.setFactory(factory)
}

// AddObjectSource adds the object source to the registry.
func (r *Registry) AddObjectSource(source IObjectSource) error {
	return r.entries.addSource(source)
}

// RemoveObjectSource removes the object source from the registry.
func (r *Registry) RemoveObjectSource(source IObjectSource) {
	if source == nil {
		log.Warn().Msg("registry: source is nil")
		return
	}
	r.entries.removeEntry(source.ObjectId())
}

// GetObjectSource returns the object source by name.
func (r *Registry) GetObjectSource(objectId string) IObjectSource {
	return r.entries.getSource(objectId)
}

// Checks if the object is registered.
func (r *Registry) IsRegistered(objectId string) bool {
	return r.entries.hasEntry(objectId)
}

// GetRemoteNode returns the node that is linked to the object.
func (r *Registry) GetRemoteNodes(objectId string) []*Node {
	return r.entries.getNodes(objectId)
}

// AttachRemoteNode attaches the node to the registry.
func (r *Registry) AttachRemoteNode(node *Node) {
}

// DetachRemoteNode removes the link between the object and the node.
func (r *Registry) DetachRemoteNode(node *Node) {
	r.entries.purgeNode(node)
}

// LinkRemoteNode adds a link between the object and the node.
func (r *Registry) LinkRemoteNode(objectId string, node *Node) {
	r.entries.addNode(objectId, node)
}

// UnlinkRemoteNode removes the link between the object and the node.
func (r *Registry) UnlinkRemoteNode(objectId string, node *Node) {
	r.entries.removeNode(objectId, node)
}

// NotifyPropertyChange notifies the property change to the nodes.
func (r *Registry) NotifyPropertyChange(objectId string, kwargs core.KWArgs) {
	log.Debug().Msgf("registry: notify property change %s", objectId)
	nodes := r.entries.getNodes(objectId)
	for _, n := range nodes {
		for name, value := range kwargs {
			propertyId := core.MakeSymbolId(objectId, name)
			n.SendPropertyChange(propertyId, value)
		}
	}
}

// NotifySignal notifies the signal to the nodes that are linked to the object.
func (r *Registry) NotifySignal(objectId string, name string, args core.Args) {
	log.Debug().Msgf("registry: notify signal %s.%s", objectId, name)
	signalId := core.MakeSymbolId(objectId, name)
	nodes := r.entries.getNodes(objectId)
	for _, n := range nodes {
		n.SendSignal(signalId, args)
	}
}
