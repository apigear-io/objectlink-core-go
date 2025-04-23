package cli

import (
	"context"

	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/ws"
)

var registry *client.Registry
var conn *ws.Connection
var node *client.Node
var ctx context.Context
var commands []Command

func init() {
	registry = client.NewRegistry()
	ctx = context.Background()
	commands = []Command{
		cmdAdd,
		cmdConnect,
		cmdDisconnect,
		cmdHelp,
		cmdInfo,
		cmdInvoke,
		cmdSignal,
		cmdLink,
		cmdUnlink,
		cmdSet,
		cmdQuit,
		cmdRemove,
		cmdGet,
		cmdServe,
	}
}
