package core

import (
	"fmt"
	"strings"
)

// Resource: <ObjectId>/<Member>
// ObjectId: <Module>.<Object>
// Resource: <Module>.<Object>/<Member>
type Resource string

func (n Resource) String() string {
	return string(n)
}

func (n Resource) ObjectId() string {
	return strings.Split((string(n)), "/")[0]
}

func (n Resource) Member() string {
	parts := strings.Split(string(n), "/")
	return parts[len(parts)-1]
}

func (n Resource) HasMember() bool {
	return strings.Contains(string(n), "/")
}

func (n Resource) IsValid() bool {
	return len(n) > 0
}

func CreateResource(objectId, member string) Resource {
	return Resource(fmt.Sprintf("%s/%s", objectId, member))
}
