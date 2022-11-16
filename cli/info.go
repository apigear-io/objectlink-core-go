package cli

import "fmt"

var info = Command{
	Usage: "info",
	Names: []string{"info"},
	Exec: func(args []string) error {
		if conn == nil {
			fmt.Printf("%s: not connected to server\n", conn.Id())
		} else {
			fmt.Printf("%s: connected to %s\n", conn.Id(), conn.Url())
		}
		if node == nil {
			fmt.Printf("not node instantiated\n")
		} else {
			fmt.Printf("node %s\n", node.Id())
			fmt.Println("registry information:")
			for _, id := range node.Registry.ObjectIds() {
				node := node.Registry.Node(id)
				if node != nil {
					fmt.Printf("  %s linked to %s\n", id, node.Id())
				} else {
					fmt.Printf("  %s not linked\n", id)
				}
			}
		}
		return nil
	},
	Help: "show status information",
}
