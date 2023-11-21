package cli

import (
	"fmt"

	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/ws"
)

var connect = Command{
	Usage: "connect <url>",
	Names: []string{"c", "con", "connect"},
	Exec: func(args []string) error {
		url := "ws://localhost:4333/ws"
		if registry == nil {
			return fmt.Errorf("no registry")
		}
		if len(args) == 2 {
			url = args[1]
		}
		c, err := ws.Dial(ctx, url)
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
