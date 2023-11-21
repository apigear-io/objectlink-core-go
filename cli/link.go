package cli

import (
	"fmt"

	"github.com/apigear-io/objectlink-core-go/log"
)

var link = Command{
	Usage: "link <objectid>",
	Names: []string{"l", "lnk", "link"},
	Exec: func(args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("no object source")
		}
		objectId := args[1]
		if registry == nil {
			return fmt.Errorf("no registry")
		}
		if node == nil {
			return fmt.Errorf("no client node")
		}
		if registry.ObjectSink(objectId) == nil {
			log.Info().Msgf("register sink for %s", objectId)
			sink := &MockSink{objectId: objectId}
			err := registry.AddObjectSink(sink)
			if err != nil {
				return err
			}
		}
		// we only have one client node
		log.Info().Msgf("link %s", objectId)
		node.LinkRemoteNode(objectId)
		return nil

	},
	Help: "connect to server",
}
