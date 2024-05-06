package cli

import (
	"fmt"
	"strings"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"gopkg.in/yaml.v3"
)

var invoke = Command{
	Usage: "invoke <methodId> [<arg>]",
	Names: []string{"i", "inv", "invoke"},
	Exec: func(args []string) error {
		if conn == nil {
			return fmt.Errorf("not connected")
		}
		if node == nil {
			return fmt.Errorf("no node")
		}
		if len(args) < 2 {
			return fmt.Errorf("missing method")
		}
		method := args[1]
		var params core.Args
		if len(args) > 2 {
			r := strings.NewReader(strings.Join(args[2:], " "))
			err := yaml.NewDecoder(r).Decode(&params)
			if err != nil {
				return err
			}
		}
		log.Info().Msgf("invoke %s %#v", method, params)
		node.InvokeRemote(method, params, func(arg client.InvokeReplyArg) {
			log.Info().Msgf("invoke reply %s", arg.Identifier)
			data, err := json.MarshalIndent(arg.Value, "", "  ")
			if err != nil {
				log.Error().Err(err).Msg("error marshalling value")
				return
			}
			fmt.Println(string(data))
		})
		return nil
	},
	Help: "invoke a method on the remote node",
}
