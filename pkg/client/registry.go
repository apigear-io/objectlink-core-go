package client

import "olink/pkg/core"

var registry *ClientRegistry

type ClientRegistry struct {
	entries map[string]*SinkToClientEntry
}

func NewClientRegistry() *ClientRegistry {
	return &ClientRegistry{
		entries: make(map[string]*SinkToClientEntry),
	}
}

func GetRegistry() *ClientRegistry {
	if registry == nil {
		registry = NewClientRegistry()
	}
	return registry
}

// attach client node to registry
func (registry *ClientRegistry) AttachClientNode(node *ClientNode) {
}

// detach client node from registry
func (registry *ClientRegistry) DetachClientNode(node *ClientNode) {
	for _, v := range registry.entries {
		if v.node == node {
			v.node = nil
		}
	}
}

func (r *ClientRegistry) LinkClientNode(name string, node *ClientNode) {
	r.Entry(name).node = node
}

func (r *ClientRegistry) UnlinkClientNode(name string) {
	r.Entry(name).node = nil
}

func (r *ClientRegistry) AddObjectSink(sink IObjectSink) {
	r.Entry(sink.ObjectName()).sink = sink
}

// remove object sink from registry
func (registry *ClientRegistry) RemoveObjectSink(sink IObjectSink) {
	name := sink.ObjectName()
	registry.RemoveEntry(name)
}

// get object sink by name
func (r *ClientRegistry) GetObjectSink(name string) IObjectSink {
	return r.Entry(name).sink
}

func (r *ClientRegistry) GetClientNode(name string) *ClientNode {
	return r.Entry(name).node
}

func (r *ClientRegistry) Entry(name string) *SinkToClientEntry {
	resource := core.ResourceFromName(name)
	if r.entries[resource] == nil {
		r.entries[resource] = &SinkToClientEntry{
			node: nil,
			sink: nil,
		}
	}
	return r.entries[resource]
}

func (r *ClientRegistry) RemoveEntry(name string) {
	resource := core.ResourceFromName(name)
	delete(r.entries, resource)
}
