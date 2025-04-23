package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/client"
	"github.com/apigear-io/objectlink-core-go/olink/core"
)

var cmdSignal = Command{
	Usage: "signal <signalId> [<arg>]",
	Names: []string{"sig", "signal"},
	Exec: func(args []string) error {
		if conn == nil {
			return fmt.Errorf("not connected")
		}
		if node == nil {
			return fmt.Errorf("no node")
		}
		if len(args) < 2 {
			return fmt.Errorf("missing signal id")
		}
		signalId := args[1]
		var params core.Args
		if len(args) > 2 {
			r := strings.NewReader(strings.Join(args[2:], " "))
			err := json.NewDecoder(r).Decode(&params)
			if err != nil {
				return err
			}
		}
		log.Info().Msgf("signal %s %#v", signalId, params)
		symbol, member := core.SymbolIdToParts(signalId)
		methodId := fmt.Sprintf("%s/$signal.%s", symbol, member)
		node.InvokeRemote(methodId, params, func(arg client.InvokeReplyArg) {
		})
		return nil
	},
	Help: "emit a signal on the remote node",
}
