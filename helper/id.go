package helper

import (
	"strconv"
	"sync/atomic"
)

func MakeIdGenerator(prefix string) func() string {
	var id atomic.Int32
	return func() string {
		next := id.Add(1)
		return prefix + strconv.Itoa(int(next))
	}
}
