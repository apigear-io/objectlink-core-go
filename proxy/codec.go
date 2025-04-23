package proxy

import (
	"bytes"
	"encoding/json"
)

type Codec interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

func NewCodec(codec string) Codec {
	switch codec {
	case "json":
		return &JsonCodec{}
	default:
		return &JsonCodec{}
	}
}

type JsonCodec struct {
}

func (c *JsonCodec) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *JsonCodec) Decode(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	// decoder.UseNumber()
	return decoder.Decode(&v)
}
