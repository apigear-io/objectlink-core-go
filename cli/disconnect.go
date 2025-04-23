package cli

import "fmt"

var cmdDisconnect = Command{
	Usage: "disconnect",
	Names: []string{"d", "disconnect"},
	Exec: func(args []string) error {
		if conn == nil {
			return fmt.Errorf("not connected")
		}
		conn.Close()
		conn = nil
		if node != nil {
			node.Close()
			node = nil
		}
		fmt.Println("disconnected")
		return nil
	},
	Help: "disconnect from server",
}
