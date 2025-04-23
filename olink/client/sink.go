package client

import (
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

type IObjectSink interface {
	ObjectId() string
	HandleSignal(signalId string, args core.Args)
	HandlePropertyChange(propertyId string, value core.Any)
	HandleInit(objectId string, props core.KWArgs, node *Node)
	HandleRelease()
}
