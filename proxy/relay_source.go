package proxy

import (
	"time"

	"github.com/apigear-io/objectlink-core-go/log"
	"github.com/apigear-io/objectlink-core-go/olink/core"
	"github.com/apigear-io/objectlink-core-go/olink/remote"
	"github.com/nats-io/nats.go"
)

const (
	DefaultTimeout = time.Second * 5
)

type RelaySource struct {
	id    string
	conn  *nats.Conn
	codec Codec
}

var _ remote.IObjectSource = (*RelaySource)(nil)

func NewRelaySource(id string, conn *nats.Conn, codec Codec) *RelaySource {
	log.Info().Msgf("relay source: new: %s", id)
	return &RelaySource{
		id:    id,
		conn:  conn,
		codec: codec,
	}
}

func (s *RelaySource) ObjectId() string {
	log.Info().Msgf("relay source: object id: %s", s.id)
	return s.id
}

func (s *RelaySource) Invoke(methodId string, args core.Args) (core.Any, error) {
	log.Info().Msgf("relay source: invoke: %s.%s", s.id, methodId)
	subj := s.id + "." + methodId
	data, err := s.codec.Encode(args)
	if err != nil {
		return nil, err
	}
	msg, err := s.conn.Request(subj, data, DefaultTimeout)
	if err != nil {
		return nil, err
	}
	var result core.Any
	err = s.codec.Decode(msg.Data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *RelaySource) SetProperty(propertyId string, value core.Any) error {
	log.Info().Msgf("relay source: set property: %s.%s", s.id, propertyId)
	subj := s.id + "." + propertyId
	data, err := s.codec.Encode(value)
	if err != nil {
		return err
	}
	_, err = s.conn.Request(subj, data, DefaultTimeout)
	if err != nil {
		return err
	}
	return nil

}
func (s *RelaySource) Linked(objectId string, node *remote.Node) error {
	log.Info().Msgf("relay source: linked: %s", objectId)
	return nil
}
func (s *RelaySource) CollectProperties() (core.KWArgs, error) {
	log.Info().Msgf("relay source: collect properties: %s", s.id)
	return core.KWArgs{}, nil
}
