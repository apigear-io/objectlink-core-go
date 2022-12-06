package core

func AsInt(v Any) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}

func AsArgs(v Any) Args {
	switch v := v.(type) {
	case []any:
		return v
	default:
		return Args{}
	}
}

func AsString(v Any) string {
	switch v := v.(type) {
	case string:
		return v
	default:
		return ""
	}
}

func AsMsgType(v Any) MsgType {
	switch v := v.(type) {
	case MsgType:
		return v
	case int:
		return MsgType(v)
	case float64:
		return MsgType(v)
	default:
		return 0
	}
}

func AsAny(v Any) Any {
	return v
}

func AsProps(v Any) Props {
	switch v := v.(type) {
	case map[string]any:
		return v
	default:
		return Props{}
	}
}

func AsResource(v Any) Resource {
	switch v := v.(type) {
	case Resource:
		return v
	case string:
		return Resource(v)
	default:
		return Resource("")
	}
}
