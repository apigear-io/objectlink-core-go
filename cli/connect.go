package cli

import (
	"fmt"
	"olink/pkg/client"
	"olink/pkg/ws"
)

var connect = Command{
	Usage: "connect <url>",
	Names: []string{"c", "connect"},
	Exec: func(args []string) error {
		url := "ws://localhost:8080/ws"
		if len(args) == 2 {
			url = args[1]
		}
		c, err := ws.Dial(url)
		if err != nil {
			return err
		}
		conn = c
		node = client.NewNode(registry)
		node.SetOutput(conn)
		conn.SetOutput(node)
		fmt.Printf("connection %s connected to %s using node %s\n", conn.Id(), url, node.Id())
		return nil
	},
	Help: "connect to server",
}
