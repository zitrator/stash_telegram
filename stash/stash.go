package stash

// todo: transaction log
// todo: encrypt data on disk

import (
	"sync"
)

// Stash in-memory storage
type Stash struct {
	sync.RWMutex
	m map[string]interface{}
}

func (s *Stash) Put(key string, doc interface{}) error {
	s.m[key] = doc
	return nil
}

func (s *Stash) Get(key string) (interface{}, error) {
	return s.m[key], nil
}
