package cli

import (
	"fmt"
	"io"
)

var quit = Command{
	Usage: "quit",
	Names: []string{"q", "quit"},
	Exec: func(args []string) error {
		if conn != nil {
			conn.Close()
		}
		fmt.Printf("bye\n")
		return io.EOF
	},
	Help: "quit",
}
