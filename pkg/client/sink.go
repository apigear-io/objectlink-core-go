package client

import (
	"olink/pkg/core"
)

type IObjectSink interface {
	ObjectId() string
	OnSignal(signalId string, args core.Args)
	OnPropertyChange(propertyId string, value core.Any)
	OnInit(objectId string, props core.KWArgs, node *Node)
	OnRelease()
}
