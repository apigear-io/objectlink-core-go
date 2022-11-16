package cli

import (
	"fmt"
)

var add = Command{
	Usage: "add <objectid>",
	Names: []string{"add", "a"},
	Exec: func(args []string) error {
		// add new object sink
		if len(args) < 2 {
			return fmt.Errorf("missing object sink")
		}
		objectId := args[1]
		if registry.ObjectSink(objectId) == nil {
			fmt.Printf("register new sink for object %s\n", objectId)
			sink := &MockSink{objectId: objectId}
			registry.AddObjectSink(sink)
		} else {
			fmt.Printf("sink for object %s already registered\n", objectId)
		}
		return nil
	},
	Help: "add new object sink",
}
