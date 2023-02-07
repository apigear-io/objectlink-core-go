package cli

import "fmt"

var link = Command{
	Usage: "link <objectid>",
	Names: []string{"l", "lnk", "link"},
	Exec: func(args []string) error {
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
		if registry.ObjectSink(objectId) == nil {
			fmt.Printf("register new sink for object %s\n", objectId)
			sink := &MockSink{objectId: objectId}
			registry.AddObjectSink(sink)
		}
		// we only have one client node
		node.LinkRemoteNode(objectId)
		return nil

	},
	Help: "connect to server",
}
