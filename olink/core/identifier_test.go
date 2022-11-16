package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToIdentifier(t *testing.T) {
	assert.Equal(t, "demo.Counter", ToObjectId("demo.Counter/member/member2"))
	assert.Equal(t, "demo.Counter", ToObjectId("demo.Counter/member"))
	assert.Equal(t, "demo.Counter", ToObjectId("demo.Counter"))
	assert.Equal(t, "", ToObjectId(""))
}

func TestToMember(t *testing.T) {
	assert.Equal(t, "", ToMember("demo.Counter/member/member2"))
	assert.Equal(t, "member", ToMember("demo.Counter/member"))
	assert.Equal(t, "", ToMember("demo.Counter"))
	assert.Equal(t, "", ToMember(""))
}

func TestToParts(t *testing.T) {
	part0, part1 := ToParts("demo.Counter/member/member2")
	assert.Equal(t, "demo.Counter", part0)
	assert.Equal(t, "", part1)
	part0, part1 = ToParts("demo.Counter/member")
	assert.Equal(t, "demo.Counter", part0)
	assert.Equal(t, "member", part1)
	part0, part1 = ToParts("demo.Counter")
	assert.Equal(t, "demo.Counter", part0)
	assert.Equal(t, "", part1)
	part0, part1 = ToParts("")
	assert.Equal(t, "", part0)
	assert.Equal(t, "", part1)
}
