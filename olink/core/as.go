package core

import (
	"encoding/json"
	"strconv"

	"github.com/apigear-io/objectlink-core-go/log"
)

func AsBool(v any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	default:
		log.Warn().Msgf("as bool unknown type %#v %T", v, v)
		return false
	}
}

func AsFloat(v any) float64 {
	switch v := v.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case json.Number:
		i, err := v.Float64()
		if err != nil {
			log.Warn().Msgf("error: %v", err)
			return 0
		}
		return i
	default:
		log.Warn().Msgf("as float unknown type %#v %T", v, v)
		return 0
	}
}

func AsInt(v any) int64 {
	switch v := v.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			log.Warn().Msgf("error: %v", err)
			return 0
		}
		return i
	default:
		log.Warn().Msgf("as int unknown type %#v %T", v, v)
		return 0
	}
}

func AsArgs(v any) Args {
	if v == nil {
		return []any{}
	}
	switch v := v.(type) {
	case []any:
		return v
	default:
		log.Warn().Msgf("as args unknown type %#v %T", v, v)
		return nil
	}
}

func AsString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		log.Warn().Msgf("as string unknown type %#v %T", v, v)
		return ""
	}
}

func AsMsgType(v any) MsgType {
	switch v := v.(type) {
	case MsgType:
		return v
	case int:
		return MsgType(v)
	case int64:
		return MsgType(v)
	case float64:
		return MsgType(v)
	case string:
		return MsgTypeFromString(v)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return MsgType(0)
		}
		return MsgType(i)
	default:
		log.Warn().Msgf("as msgtype unknown type %#v %T", v, v)
		return 0
	}
}

func AsAny(v any) any {
	return v
}

func AsProps(v any) KWArgs {
	if v == nil {
		return KWArgs{}
	}
	switch v := v.(type) {
	case map[string]any:
		return v
	default:
		log.Warn().Msgf("as props unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayBool(v any) []bool {
	switch v := v.(type) {
	case []bool:
		return v
	default:
		log.Warn().Msgf("as array bool unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayInt(v any) []int64 {
	if v == nil {
		return []int64{}
	}
	switch v := v.(type) {
	case []int64:
		return v
	case []any:
		r := make([]int64, len(v))
		for i, v := range v {
			r[i] = AsInt(v)
		}
		return r
	default:
		log.Warn().Msgf("as array int unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayFloat(v any) []float64 {
	if v == nil {
		return []float64{}
	}
	switch v := v.(type) {
	case []float64:
		return v
	case []any:
		r := make([]float64, len(v))
		for i, v := range v {
			r[i] = AsFloat(v)
		}
		return r
	default:
		log.Warn().Msgf("as array float unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayString(v any) []string {
	if v == nil {
		return []string{}
	}
	switch v := v.(type) {
	case []string:
		return v
	case []any:
		r := make([]string, len(v))
		for i, v := range v {
			r[i] = AsString(v)
		}
		return r
	default:
		log.Warn().Msgf("as array string unknown type %#v %T", v, v)
		return nil
	}
}

func AsStruct(v any) KWArgs {
	if v == nil {
		return KWArgs{}
	}
	switch v := v.(type) {
	case map[string]any:
		return v
	case KWArgs:
		return v
	default:
		log.Warn().Msgf("as struct unknown type %#v %T", v, v)
		return nil
	}
}

func AsEnum(v any) []any {
	if v == nil {
		return []any{}
	}
	switch v := v.(type) {
	case []any:
		return v
	default:
		log.Warn().Msgf("as enum unknown type %#v %T", v, v)
		return nil
	}
}

func AsInterface(v any) interface{} {
	return v
}

func AsArrayStruct(v any) []KWArgs {
	if v == nil {
		return []KWArgs{}
	}
	switch v := v.(type) {
	case []KWArgs:
		return v
	case []map[string]any:
		r := make([]KWArgs, len(v))
		for i, v := range v {
			r[i] = v
		}
		return r
	default:
		log.Warn().Msgf("as array struct unknown type %#v %T", v, v)
		return []KWArgs{}
	}
}

func AsArrayEnum(v any) [][]any {
	switch v := v.(type) {
	case [][]any:
		return v
	default:
		log.Warn().Msgf("as array enum unknown type %#v %T", v, v)
		return [][]any{}
	}
}

func AsArrayInterface(v any) []interface{} {
	switch v := v.(type) {
	case []interface{}:
		return v
	default:
		log.Warn().Msgf("as array interface unknown type %#v %T", v, v)
		return []interface{}{}
	}
}
