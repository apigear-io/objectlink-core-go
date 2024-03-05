package proxy

import (
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/remote"
)

// ObjectLink Proxy does proxy the ObjectLink messages from a sink to a source registry and vice versa.
// It uses the olink client node as a sink and the olink remote node as a source.

type Proxy struct {
	codec Codec
	// source registry
	source *remote.Registry
	// sink registry
	sink *client.Registry
}

// NewProxy creates a new proxy instance.
func NewProxy() *Proxy {
	return &Proxy{
		codec:  NewCodec("json"),
		source: remote.NewRegistry(),
		sink:   client.NewRegistry(),
	}
}

func (p *Proxy) Codec() Codec {
	return p.codec
}

func (p *Proxy) SourceRegistry() *remote.Registry {
	return p.source
}

func (p *Proxy) SinkRegistry() *client.Registry {
	return p.sink
}

func (p *Proxy) CreateClientNode() *client.Node {
	return client.NewNode(p.sink)
}

func (p *Proxy) CreateRemoteNode() *remote.Node {
	return remote.NewNode(p.source)
}
