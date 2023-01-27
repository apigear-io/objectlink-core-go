package core

import (
	"fmt"
	"strings"
)

// Identifier: <object-id>/<member>
// ObjectId: <module-name>.<object-name>
// Identifier: <module-name>.<object-name>/<member>

func SymbolIdToObjectId(id string) string {
	return strings.Split(id, "/")[0]
}

func SymbolIdToMember(id string) string {
	parts := strings.Split(id, "/")
	switch len(parts) {
	case 1:
		return ""
	case 2:
		return parts[1]
	default:
		return ""
	}
}

func SymbolIdToParts(id string) (string, string) {
	parts := strings.Split(id, "/")
	switch len(parts) {
	case 1:
		return parts[0], ""
	case 2:
		return parts[0], parts[1]
	default:
		return parts[0], ""
	}
}

func MakeSymbolId(id string, member string) string {
	return fmt.Sprintf("%s/%s", id, member)
}
