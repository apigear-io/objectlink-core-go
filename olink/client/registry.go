package client

type SinkToClientEntry struct {
	sink IObjectSink
	node *Node
}

type Registry struct {
	entries map[string]*SinkToClientEntry
}

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]*SinkToClientEntry),
	}
}

// attach client node to registry
func (registry *Registry) AttachClientNode(node *Node) {
}

// detach client node from registry
func (registry *Registry) DetachClientNode(node *Node) {
	for _, v := range registry.entries {
		if v.node == node {
			v.node = nil
		}
	}
}

func (r *Registry) LinkClientNode(objectId string, node *Node) {
	r.Entry(objectId).node = node
}

func (r *Registry) UnlinkClientNode(objectId string) {
	r.Entry(objectId).node = nil
}

func (r *Registry) AddObjectSink(sink IObjectSink) {
	r.Entry(sink.ObjectId()).sink = sink
}

// remove object sink from registry
func (registry *Registry) RemoveObjectSink(sink IObjectSink) {
	objectId := sink.ObjectId()
	registry.RemoveEntry(objectId)
}

// get object sink by name
func (r *Registry) GetObjectSink(objectId string) IObjectSink {
	return r.Entry(objectId).sink
}

func (r *Registry) GetClientNode(objectId string) *Node {
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
