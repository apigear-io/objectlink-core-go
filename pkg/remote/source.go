package remote

import "olink/pkg/core"

type IObjectSource interface {
	ObjectId() string
	Invoke(res core.Resource, args core.Args)
	SetProperty(res core.Resource, value core.Any)
	Linked(objectId string, node *Node)
	CollectProperties() core.Props
}
