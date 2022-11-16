package cli

import (
	"olink/pkg/client"
	"olink/pkg/ws"
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
