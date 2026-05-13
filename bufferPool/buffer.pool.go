package bufferpool

import (
	"sync"
)

const MAX_SIZE = 1024

var Bufferpool = sync.Pool{
	New: func() any {
		b := make([]string, 0, MAX_SIZE)
		return &b
	},
}
