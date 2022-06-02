package stash

import (
	"sync"
)

// Stash in-memory storage
// todo: transaction log
// todo: encrypt data on disk
type Stash struct {
	sync.RWMutex
	m map[string]string
}
