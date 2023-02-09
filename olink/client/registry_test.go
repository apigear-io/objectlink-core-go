package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextRegistryId(t *testing.T) {
	t.Parallel()
	clearRegistryId()
	id := nextRegistryId()
	assert.Equal(t, "r1", id)
	id = nextRegistryId()
	assert.Equal(t, "r2", id)
	id = nextRegistryId()
	assert.Equal(t, "r3", id)
}

func TestSetSinkFactory(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	r.SetSinkFactory(nil)

	s := r.ObjectSink("demo.Counter")
	assert.Nil(t, s)

	factory := func(objectId string) IObjectSink {
		return NewMockSink(objectId)
	}
	r.SetSinkFactory(factory)

	s = r.ObjectSink("demo.Counter")
	assert.NotNil(t, s)
	assert.Equal(t, "demo.Counter", s.ObjectId())

	s2 := r.ObjectSink("demo.Counter")
	assert.Equal(t, s, s2)
}

func TestAddObjectSink(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSink("demo.Counter")
	err := r.AddObjectSink(s)
	assert.Nil(t, err)
	assert.Equal(t, s, r.ObjectSink("demo.Counter"))

	s2 := NewMockSink("demo.Counter")
	err = r.AddObjectSink(s2)
	assert.NotNil(t, err)
	assert.Equal(t, s, r.ObjectSink("demo.Counter"))
}

func TestIsRegistered(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	assert.False(t, r.IsRegistered("demo.Counter"))
	s := NewMockSink("demo.Counter")
	r.AddObjectSink(s)
	assert.True(t, r.IsRegistered("demo.Counter"))
}

func TestRemoveObjectSink(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSink("demo.Counter")
	r.AddObjectSink(s)
	assert.True(t, r.IsRegistered("demo.Counter"))
	r.RemoveObjectSink("demo.Counter")
	assert.False(t, r.IsRegistered("demo.Counter"))
}

func TestGetObjectSink(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSink("demo.Counter")
	r.AddObjectSink(s)
	assert.Equal(t, s, r.ObjectSink("demo.Counter"))
}

func TestGetClientNode(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	assert.Nil(t, r.GetClientNode("demo.Counter"))
	n1 := NewNode(r)
	s1 := NewMockSink("demo.Counter")
	r.AddObjectSink(s1)
	r.LinkClientNode(s1.ObjectId(), n1)
	assert.Equal(t, n1, r.GetClientNode(s1.ObjectId()))
	r.LinkClientNode(s1.ObjectId(), n1)
	assert.Equal(t, n1, r.GetClientNode(s1.ObjectId()))
}

func TestDetachClientNode(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	n1 := NewNode(r)
	s1 := NewMockSink("demo.Counter")
	r.AddObjectSink(s1)
	r.LinkClientNode(s1.ObjectId(), n1)
	assert.Equal(t, n1, r.GetClientNode(s1.ObjectId()))
	r.DetachClientNode(n1)
	assert.Nil(t, r.GetClientNode(s1.ObjectId()))
	r.DetachClientNode(n1)
	assert.Nil(t, r.GetClientNode(s1.ObjectId()))
}
