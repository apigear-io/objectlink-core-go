package cli

import "fmt"

var remove = Command{
	Usage: "remove <objectId>",
	Names: []string{"r", "remove"},
	Exec: func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("missing object sink")
		}
		objectId := args[1]
		if registry == nil {
			return fmt.Errorf("no registry")
		}
		if registry.ObjectSink(objectId) == nil {
			fmt.Printf("object sink %s not found\n", objectId)
		} else {
			fmt.Printf("remove object sink %s\n", objectId)
			registry.RemoveObjectSink(objectId)
		}
		return nil
	},
	Help: "remove object sink",
}
