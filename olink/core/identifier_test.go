package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToIdentifier(t *testing.T) {
	assert.Equal(t, "demo.Counter", SymbolIdToObjectId("demo.Counter/member/member2"))
	assert.Equal(t, "demo.Counter", SymbolIdToObjectId("demo.Counter/member"))
	assert.Equal(t, "demo.Counter", SymbolIdToObjectId("demo.Counter"))
	assert.Equal(t, "", SymbolIdToObjectId(""))
}

func TestToMember(t *testing.T) {
	assert.Equal(t, "", SymbolIdToMember("demo.Counter/member/member2"))
	assert.Equal(t, "member", SymbolIdToMember("demo.Counter/member"))
	assert.Equal(t, "", SymbolIdToMember("demo.Counter"))
	assert.Equal(t, "", SymbolIdToMember(""))
}

func TestToParts(t *testing.T) {
	part0, part1 := SymbolIdToParts("demo.Counter/member/member2")
	assert.Equal(t, "demo.Counter", part0)
	assert.Equal(t, "", part1)
	part0, part1 = SymbolIdToParts("demo.Counter/member")
	assert.Equal(t, "demo.Counter", part0)
	assert.Equal(t, "member", part1)
	part0, part1 = SymbolIdToParts("demo.Counter")
	assert.Equal(t, "demo.Counter", part0)
	assert.Equal(t, "", part1)
	part0, part1 = SymbolIdToParts("")
	assert.Equal(t, "", part0)
	assert.Equal(t, "", part1)
}
