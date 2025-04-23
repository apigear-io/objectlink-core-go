package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/remote"
	"github.com/apigear-io/objectlink-core-go/olink/ws"
)

type GenericSource struct {
	objectId   string
	properties core.KWArgs
	node       *remote.Node
}

var _ remote.IObjectSource = (*GenericSource)(nil)

func NewGenericSource(objectId string) *GenericSource {
	log.Info().Str("objectId", objectId).Msg("create new source")
	return &GenericSource{
		objectId:   objectId,
		properties: make(core.KWArgs),
	}
}

func (s *GenericSource) ObjectId() string {
	log.Info().Str("objectId", s.objectId).Msg("get objectId")
	return s.objectId
}

func (s *GenericSource) Invoke(methodId string, args core.Args) (core.Any, error) {
	log.Info().Str("objectId", s.objectId).Msgf("invoke method %s with args %v", methodId, args)
	switch {
	case strings.HasPrefix(methodId, "$get"):
		switch len(args) {
		case 0:
			return s.properties, nil
		case 1:
			if value, ok := s.properties[methodId]; ok {
				return value, nil
			}
		}
	case strings.HasPrefix(methodId, "$signal"):
		signal := methodId[len("$signal."):]
		if s.node == nil {
			log.Error().Str("objectId", s.objectId).Msg("node is not set")
			return nil, fmt.Errorf("node is not set")
		}
		log.Info().Str("objectId", s.objectId).Msgf("signal %s with args %v", signal, args)
		s.node.Registry().NotifySignal(s.objectId, signal, args)
	}
	return nil, nil
}
func (s *GenericSource) SetProperty(propertyId string, value core.Any) error {
	log.Info().Str("objectId", s.objectId).Msgf("set property %s to %v", propertyId, value)
	s.properties[propertyId] = value
	if s.node != nil {
		s.node.Registry().NotifyPropertyChange(s.objectId, core.KWArgs{
			propertyId: value,
		})
	}
	return nil
}
func (s *GenericSource) Linked(objectId string, node *remote.Node) error {
	log.Info().Str("objectId", s.objectId).Msgf("linked to %s", objectId)
	if s.objectId != objectId {
		err := fmt.Errorf("objectId mismatch %s != %s", s.objectId, objectId)
		log.Error().Err(err).Msg("source linked error")
		return err
	}
	s.node = node
	return nil
}
func (s *GenericSource) CollectProperties() (core.KWArgs, error) {
	return s.properties, nil
}

func GenericSourceFactory(objectId string) remote.IObjectSource {
	return NewGenericSource(objectId)
}

func RunHub(addr string) {
	registry := remote.NewRegistry()
	registry.SetSourceFactory(GenericSourceFactory)
	hub := ws.NewHub(ctx, registry)
	server := &http.Server{
		Addr: addr,
	}
	http.Handle("/ws", hub)

	go func() {
		log.Info().Msgf("objectlink server listening on ws://%s/ws", addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("failed to start web socket server")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	ctx := context.Background()
	err := server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown web socket server")
	}
	log.Info().Msg("web socket server shutdown")
}

var cmdServe = Command{
	Usage: "serve <addr>",
	Names: []string{"s", "serve"},
	Exec: func(args []string) error {
		addr := "localhost:5555"
		if len(args) > 1 {
			addr = args[1]
		}
		RunHub(addr)
		return nil
	},
	Help: "start an objectlink server",
}
