package client

import "olink/pkg/core"

type IObjectSink interface {
	ObjectId() string
	OnSignal(res core.Resource, args core.Args)
	OnPropertyChange(res core.Resource, value core.Any)
	OnInit(objectId string, props core.KWArgs, node *Node)
	OnRelease()
}
