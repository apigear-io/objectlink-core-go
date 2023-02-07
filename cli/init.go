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
var cancel context.CancelFunc
var commands []Command

func init() {
	registry = client.NewRegistry()
	ctx, cancel = context.WithCancel(context.Background())
	commands = []Command{
		add,
		connect,
		disconnect,
		help,
		info,
		invoke,
		link,
		quit,
		remove,
		set,
		unlink,
	}
}
