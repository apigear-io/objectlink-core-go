package remote

import (
	"testing"

	"github.com/apigear-io/objectlink-core-go/helper"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/stretchr/testify/require"
)

func TestNextRegistryId(t *testing.T) {
	nextRegistryId := helper.MakeIdGenerator("r")
	id := nextRegistryId()
	require.Equal(t, "r1", id)
	id = nextRegistryId()
	require.Equal(t, "r2", id)
	id = nextRegistryId()
	require.Equal(t, "r3", id)
}

func TestSetSourceFactory(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	r.SetSourceFactory(nil)

	s := r.GetObjectSource("demo.Counter")
	require.Nil(t, s)

	factory := func(objectId string) IObjectSource {
		return NewMockSource(objectId)
	}
	r.SetSourceFactory(factory)

	s = r.GetObjectSource("demo.Counter")
	require.NotNil(t, s)
	require.Equal(t, "demo.Counter", s.ObjectId())

	s2 := r.GetObjectSource("demo.Counter")
	require.Equal(t, s, s2)
}

func TestAddObjectSource(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	err := r.AddObjectSource(s)
	require.Nil(t, err)
	require.Equal(t, s, r.GetObjectSource("demo.Counter"))

	s2 := NewMockSource("demo.Counter")
	err = r.AddObjectSource(s2)
	require.NotNil(t, err)
	require.Equal(t, s, r.GetObjectSource("demo.Counter"))
}

func TestIsRegistered(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	require.False(t, r.IsRegistered("demo.Counter"))
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	require.True(t, r.IsRegistered("demo.Counter"))
}

func TestRemoveObjectSource(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	require.True(t, r.IsRegistered("demo.Counter"))
	r.RemoveObjectSource(s)
	require.False(t, r.IsRegistered("demo.Counter"))

	r.RemoveObjectSource(s)
	require.False(t, r.IsRegistered("demo.Counter"))
}

func TestGetObjectSource(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := r.GetObjectSource("demo.Counter")
	require.Nil(t, s)

	factory := func(objectId string) IObjectSource {
		return NewMockSource(objectId)
	}
	r.SetSourceFactory(factory)
	s = r.GetObjectSource("demo.Counter")
	require.NotNil(t, s)
	require.Equal(t, "demo.Counter", s.ObjectId())
}

func TestGetRemoteNodes(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	n := NewNode(r)
	wc := NewMockWriteCloser()
	n.SetOutput(wc)
	ns := r.GetRemoteNodes("demo.Counter")
	require.Equal(t, 0, len(ns))

	r.LinkRemoteNode("demo.Counter", n)

	ns = r.GetRemoteNodes("demo.Counter")
	require.Equal(t, 1, len(ns))
}

func TestDetachRemoteNode(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	s2 := NewMockSource("demo.Storage")
	r.AddObjectSource(s)
	r.AddObjectSource(s2)
	n := NewNode(r)
	// no nodes attached
	require.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
	require.Equal(t, 0, len(r.GetRemoteNodes(s2.ObjectId())))
	r.LinkRemoteNode(s.ObjectId(), n)
	r.LinkRemoteNode(s2.ObjectId(), n)
	// attached nodes
	require.Equal(t, 1, len(r.GetRemoteNodes(s.ObjectId())))
	require.Equal(t, 1, len(r.GetRemoteNodes(s2.ObjectId())))
	r.DetachRemoteNode(n)
	// no nodes attached
	require.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
	require.Equal(t, 0, len(r.GetRemoteNodes(s2.ObjectId())))
	r.DetachRemoteNode(n)
	require.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
	require.Equal(t, 0, len(r.GetRemoteNodes(s2.ObjectId())))
}

func TestLinkRemoteNode(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	n := NewNode(r)
	require.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
	r.LinkRemoteNode(s.ObjectId(), n)
	require.Equal(t, 1, len(r.GetRemoteNodes(s.ObjectId())))
	r.LinkRemoteNode(s.ObjectId(), n)
	require.Equal(t, 1, len(r.GetRemoteNodes(s.ObjectId())))
}

func TestUnlinkRemoteNode(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	n := NewNode(r)
	r.LinkRemoteNode(s.ObjectId(), n)
	require.Equal(t, 1, len(r.GetRemoteNodes(s.ObjectId())))
	r.UnlinkRemoteNode(s.ObjectId(), n)
	require.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
	r.UnlinkRemoteNode(s.ObjectId(), n)
	require.Equal(t, 0, len(r.GetRemoteNodes(s.ObjectId())))
}

func TestNotifyPropertyChange(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	n := NewNode(r)
	wc := NewMockWriteCloser()
	n.SetOutput(wc)
	r.LinkRemoteNode(s.ObjectId(), n)
	n.NotifyPropertyChange("demo.Counter/count", 10)
	require.Equal(t, 1, len(wc.Messages))
	msg, err := n.conv.FromData(wc.Messages[0])
	require.Nil(t, err)
	propId, value := msg.AsPropertyChange()
	require.Equal(t, "demo.Counter/count", propId)
	require.Equal(t, int64(10), core.AsInt(value))
}

func TestMultiNodePropertyChange(t *testing.T) {
	t.Parallel()
	r := NewRegistry()
	s := NewMockSource("demo.Counter")
	r.AddObjectSource(s)
	n := NewNode(r)
	wc := NewMockWriteCloser()
	n.SetOutput(wc)
	r.LinkRemoteNode(s.ObjectId(), n)
	n2 := NewNode(r)
	wc2 := NewMockWriteCloser()
	n2.SetOutput(wc2)
	r.LinkRemoteNode(s.ObjectId(), n2)
	r.NotifyPropertyChange(s.ObjectId(), core.KWArgs{"count": 10})
	require.Equal(t, 1, len(wc.Messages))
	require.Equal(t, 1, len(wc2.Messages))
	msg, err := n.conv.FromData(wc.Messages[0])
	require.Nil(t, err)
	propId, value := msg.AsPropertyChange()
	require.Equal(t, "demo.Counter/count", propId)
	require.Equal(t, int64(10), core.AsInt(value))
	msg, err = n2.conv.FromData(wc2.Messages[0])
	require.Nil(t, err)
	propId, value = msg.AsPropertyChange()
	require.Equal(t, "demo.Counter/count", propId)
	require.Equal(t, int64(10), core.AsInt(value))
}
