package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsBool(t *testing.T) {
	t.Parallel()
	assert.False(t, AsBool(nil))
	assert.False(t, AsBool(0))
	assert.False(t, AsBool(0.0))
}

func TestAsFloat(t *testing.T) {
	t.Parallel()
	assert.Equal(t, 0.0, AsFloat(nil))
	assert.Equal(t, 0.0, AsFloat(0))
	assert.Equal(t, 0.0, AsFloat(0.0))
	assert.Equal(t, 1.0, AsFloat(1))
}

func TestAsInt(t *testing.T) {
	t.Parallel()
	assert.Equal(t, int64(0), AsInt(nil))
	assert.Equal(t, int64(0), AsInt(0))
	assert.Equal(t, int64(0), AsInt(0.0))
	assert.Equal(t, int64(1), AsInt(1))
}

func TestAsArgs(t *testing.T) {
	t.Parallel()
	assert.Equal(t, Args{}, AsArgs(nil))
	assert.Equal(t, Args{}, AsArgs([]any{}))
	assert.Equal(t, Args{1, 2, 3}, AsArgs([]any{1, 2, 3}))
}

func TestAsString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "", AsString(nil))
	assert.Equal(t, "0", AsString(0))
	assert.Equal(t, "", AsString(0.0))
	assert.Equal(t, "1", AsString(1))
}

func TestAsMsgType(t *testing.T) {
	t.Parallel()
	assert.Equal(t, MsgLink, AsMsgType("link"))
	assert.Equal(t, MsgLink, AsMsgType(10))
	assert.Equal(t, MsgLink, AsMsgType(10.0))
}

func TestAsProps(t *testing.T) {
	t.Parallel()
	assert.Equal(t, KWArgs{}, AsProps(nil))
	assert.Equal(t, KWArgs{}, AsProps(map[string]any{}))
}

func TestAsArrayInt(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []int64{}, AsArrayInt(nil))
	assert.Equal(t, []int64{}, AsArrayInt([]any{}))
	assert.Equal(t, []int64{1, 2, 3}, AsArrayInt([]any{1, 2, 3}))
}

func TestArrayFloat(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []float64{}, AsArrayFloat(nil))
	assert.Equal(t, []float64{}, AsArrayFloat([]any{}))
	assert.Equal(t, []float64{1, 2, 3}, AsArrayFloat([]any{1, 2, 3}))
}

func TestAsArrayString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{}, AsArrayString(nil))
	assert.Equal(t, []string{}, AsArrayString([]any{}))
	assert.Equal(t, []string{"1", "2", "3"}, AsArrayString([]any{1, 2, 3}))
}

func TestAsStruct(t *testing.T) {
	t.Parallel()
	assert.Equal(t, KWArgs{}, AsStruct(nil))
	assert.Equal(t, KWArgs{}, AsStruct(map[string]any{}))
	assert.Equal(t, KWArgs{"a": 1, "b": 2}, AsStruct(map[string]any{"a": 1, "b": 2}))
}

func TestAsEnum(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []any{}, AsEnum(nil))
	assert.Equal(t, []any{}, AsEnum([]any{}))
	assert.Equal(t, []any{1, 2, 3}, AsEnum([]any{1, 2, 3}))
}

func TestAsArrayStruct(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []KWArgs{}, AsArrayStruct(nil))
	assert.Equal(t, []KWArgs{}, AsArrayStruct([]any{}))
	assert.Equal(t, []KWArgs{{"a": 1, "b": 2}}, AsArrayStruct([]map[string]any{{"a": 1, "b": 2}}))
}

func TestAsArrayEnum(t *testing.T) {
	t.Parallel()
	assert.Equal(t, [][]any{}, AsArrayEnum(nil))
	assert.Equal(t, [][]any{}, AsArrayEnum([]any{}))
	assert.Equal(t, [][]any{{1, 2, 3}}, AsArrayEnum([][]any{{1, 2, 3}}))
}

func TestAsArrayInterface(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []any{}, AsArrayInterface(nil))
	assert.Equal(t, []any{}, AsArrayInterface([]any{}))
	assert.Equal(t, []any{1, 2, 3}, AsArrayInterface([]any{1, 2, 3}))
}
