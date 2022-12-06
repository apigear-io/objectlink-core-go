package remote

type SourceToNodeEntry struct {
	source IObjectSource
	nodes  []*Node
}

type Registry struct {
	entries map[string]*SourceToNodeEntry
}

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]*SourceToNodeEntry),
	}
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
	return r.entry(objectId).source
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
	if r.entries[objectId] == nil {
		r.entries[objectId] = &SourceToNodeEntry{}
	}
	return r.entries[objectId]
}

func (r *Registry) removeEntry(objectId string) {
	delete(r.entries, objectId)
}
