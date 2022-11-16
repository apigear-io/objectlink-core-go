package cli

import (
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/ws"
)

var registry = client.NewRegistry()
var conn *ws.Connection
var node *client.Node

var commands []Command

func init() {
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
