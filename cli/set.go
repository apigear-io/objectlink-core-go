package cli

import "fmt"

var set = Command{
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
		name := args[1]
		value := args[2]
		node.SetRemoteProperty(name, value)
		return nil
	},
	Help: "set property value",
}
