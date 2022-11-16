package remote

import "olink/pkg/core"

type IObjectSource interface {
	ObjectId() string
	Invoke(methodId string, args core.Args) (core.Any, error)
	SetProperty(propertyId string, value core.Any) error
	Linked(objectId string, node *Node) error
	CollectProperties() (core.KWArgs, error)
}
