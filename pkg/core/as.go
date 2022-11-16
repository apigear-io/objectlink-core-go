package core

import (
	"encoding/json"
	"olink/log"
	"strconv"
)

func AsBool(v Any) bool {
	switch v := v.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return false
	}
}

func AsFloat(v Any) float64 {
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
			log.Warn().Msgf("error: %v\n", err)
			return 0
		}
		return i
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
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
			log.Warn().Msgf("error: %v\n", err)
			return 0
		}
		return i
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return 0
	}
}

func AsArgs(v Any) Args {
	switch v := v.(type) {
	case []any:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsString(v Any) string {
	switch v := v.(type) {
	case string:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return ""
	}
}

func AsMsgType(v Any) MsgType {
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
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return MsgType(0)
		}
		return MsgType(i)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return MsgType(0)
		}
		return MsgType(i)
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return 0
	}
}

func AsAny(v Any) Any {
	return v
}

func AsProps(v Any) KWArgs {
	switch v := v.(type) {
	case map[string]any:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayBool(v Any) []bool {
	switch v := v.(type) {
	case []bool:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayInt(v Any) []int64 {
	switch v := v.(type) {
	case []int64:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayFloat(v Any) []float64 {
	switch v := v.(type) {
	case []float64:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsArrayString(v Any) []string {
	switch v := v.(type) {
	case []string:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsStruct(v Any) KWArgs {
	switch v := v.(type) {
	case KWArgs:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsEnum(v Any) []any {
	switch v := v.(type) {
	case []any:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return nil
	}
}

func AsInterface(v Any) interface{} {
	return v
}

func AsArrayStruct(v Any) []KWArgs {
	switch v := v.(type) {
	case []KWArgs:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return []KWArgs{}
	}
}

func AsArrayEnum(v Any) [][]any {
	switch v := v.(type) {
	case [][]any:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return [][]any{}
	}
}

func AsArrayInterface(v Any) []interface{} {
	switch v := v.(type) {
	case []interface{}:
		return v
	default:
		log.Warn().Msgf("unknown type %#v %T", v, v)
		return []interface{}{}
	}
}
