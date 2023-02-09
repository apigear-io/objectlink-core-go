package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConverter(t *testing.T) {
	c := NewConverter(FormatJson)
	msg := MakeLinkMessage("test")
	data, err := c.ToData(msg)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	act, err := c.FromData(data)
	assert.Nil(t, err)
	assert.Equal(t, msg.AsLink(), act.AsLink())
}
