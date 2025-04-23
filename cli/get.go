package cli

import (
	"encoding/json"
	"fmt"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var cmdGet = Command{
	Usage: "get <objectId>",
	Names: []string{"g", "get"},
	Exec: func(args []string) error {
		var objectIds []string
		if conn == nil {
			return fmt.Errorf("not connected")
		}
		if node == nil {
			return fmt.Errorf("no node")
		}
		if len(args) < 2 {
			log.Info().Msg("no object id specified, getting all objects")
			objectIds = registry.ObjectIds()
		} else {
			objectIds = args[1:]
		}
		for _, objectId := range objectIds {
			if registry.ObjectSink(objectId) == nil {
				return fmt.Errorf("no sink for object %s", objectId)
			}
			method := core.MakeSymbolId(objectId, "$get")
			log.Info().Msgf("invoke %s", method)
			node.InvokeRemote(method, core.Args{}, func(arg client.InvokeReplyArg) {
				log.Info().Msgf("invoke reply %s", arg.Identifier)
				data, err := json.MarshalIndent(arg.Value, "", "  ")
				if err != nil {
					log.Error().Err(err).Msg("error marshalling value")
					return
				}
				fmt.Println(string(data))
			})
		}

		return nil
	},
	Help: "get properties of the remote object",
}
