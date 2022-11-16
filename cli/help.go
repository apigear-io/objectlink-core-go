package cli

import (
	"fmt"
	"strings"
)

var help = Command{
	Usage: "help",
	Names: []string{"h", "help", "?"},
	Exec: func(args []string) error {
		if len(args) < 2 {
			for _, cmd := range commands {
				fmt.Printf("%-32s - %s\n", cmd.Usage, cmd.Help)
			}
		} else {
			for _, cmd := range commands {
				for _, name := range cmd.Names {
					if name == args[1] {
						fmt.Printf("%-32s - %-12s - %s\n", cmd.Usage, strings.Join(cmd.Names, ","), cmd.Help)
						return nil
					}
				}
			}
			return fmt.Errorf("unknown command %s", args[1])
		}
		return nil
	},
	Help: "show commands",
}
