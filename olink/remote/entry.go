package remote

import (
	"fmt"
	"sync"
)

// sourceToNodeEntry links an object source to a list of nodes
type sourceToNodeEntry struct {
	sync.RWMutex
	source IObjectSource
	nodes  []*Node
}

// addNode adds the node to the entry
func (e *sourceToNodeEntry) addNode(node *Node) error {
	e.Lock()
	defer e.Unlock()
	for _, n := range e.nodes {
		if n == node {
			return fmt.Errorf("node already exists")
		}
	}
	e.nodes = append(e.nodes, node)
	return nil
}

// removeNode removes the node from the entry
func (e *sourceToNodeEntry) removeNode(node *Node) {
	e.Lock()
	defer e.Unlock()
	for i, n := range e.nodes {
		if n == node {
			e.nodes = append(e.nodes[:i], e.nodes[i+1:]...)
			return
		}
	}
}

// hasNode returns true if the node is linked to the source
func (e *sourceToNodeEntry) hasNode(node *Node) bool {
	e.RLock()
	defer e.RUnlock()
	for _, n := range e.nodes {
		if n == node {
			return true
		}
	}
	return false
}

// getNodes returns a copy of the list of nodes
func (e *sourceToNodeEntry) getNodes() []*Node {
	e.RLock()
	defer e.RUnlock()
	nodes := make([]*Node, len(e.nodes))
	copy(nodes, e.nodes)
	return nodes
}

// setSource sets the source
func (e *sourceToNodeEntry) setSource(source IObjectSource) {
	e.Lock()
	defer e.Unlock()
	e.source = source
}

// getSource returns the source
func (e *sourceToNodeEntry) getSource() IObjectSource {
	e.RLock()
	defer e.RUnlock()
	return e.source
}

// hasSource returns true if the source is not nil
func (e *sourceToNodeEntry) hasSource() bool {
	e.RLock()
	defer e.RUnlock()
	return e.source != nil
}
