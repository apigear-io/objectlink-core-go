package cli

import (
	"encoding/json"
	"fmt"
)

var cmdSet = Command{
	Usage: "set <propertyId> <value>",
	Names: []string{"s", "set"},
	Exec: func(args []string) error {
		// set property value
		if len(args) < 2 {
			return fmt.Errorf("no property name")
		}
		if len(args) < 3 {
			return fmt.Errorf("no property value")
		}
		if node == nil {
			return fmt.Errorf("no client node")
		}
		name := args[1]
		data := []byte(args[2])
		var value any
		err := json.Unmarshal(data, &value)
		if err != nil {
			return fmt.Errorf("failed to unmarshal property value: %w", err)
		}
		node.SetRemoteProperty(name, value)
		return nil
	},
	Help: "set property value",
}
