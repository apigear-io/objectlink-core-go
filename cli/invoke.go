package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var invoke = Command{
	Usage: "invoke <methodId> [<arg>]",
	Names: []string{"i", "invoke"},
	Exec: func(args []string) error {
		if conn == nil {
			return fmt.Errorf("not connected")
		}
		if len(args) < 2 {
			return fmt.Errorf("missing method")
		}
		method := args[1]
		var params core.Args
		if len(args) > 2 {
			r := strings.NewReader(strings.Join(args[2:], " "))
			err := json.NewDecoder(r).Decode(&params)
			if err != nil {
				return err
			}
		}
		log.Info().Msgf("invoke %s %#v", method, params)
		node.InvokeRemote(method, params, func(arg client.InvokeReplyArg) {
			fmt.Printf("%s: %v\n", arg.Identifier, arg.Value)
		})
		return nil
	},
	Help: "invoke a method on the remote node",
}
