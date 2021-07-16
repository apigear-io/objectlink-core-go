package client

import "olink/pkg/core"

type InvokeReplyArg struct {
	name  string
	value core.Any
}

type InvokeReplyFunc func(args InvokeReplyArg)

type IObjectSink interface {
	ObjectName() string
	OnSignal(name string, args core.Args)
	OnPropertyChange(name string, value core.Any)
	OnInit(name string, props core.Props, node *ClientNode)
	OnRelease()
}

type SinkToClientEntry struct {
	sink IObjectSink
	node *ClientNode
}
