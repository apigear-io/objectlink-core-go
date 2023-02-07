package cli

import "fmt"

var unlink = Command{
	Usage: "unlink <objectId>",
	Names: []string{"u", "ulnk", "unlink"},
	Exec: func(args []string) error {
		// unlink object source
		if len(args) < 2 {
			return fmt.Errorf("no object source")
		}
		objectId := args[1]
		if registry == nil {
			return fmt.Errorf("no registry")
		}
		if node == nil {
			return fmt.Errorf("no client node")
		}
		registry.RemoveObjectSink(objectId)
		node.UnlinkRemoteNode(objectId)
		return nil
	},
	Help: "disconnect from server",
}
