package stash

// todo: encrypt data on disk

import (
	"errors"
	"sync"
)

// Stash in-memory storage
type Stash struct {
	// TODO: implement io.Reader
	// TODO: implement io.Writer
	sync.RWMutex
	m map[string]interface{}
}

// NewStash constructor
func NewStash() *Stash {
	return &Stash{
		m: make(map[string]interface{}),
	}
}

// ErrorNoSuchKey predefined error
var ErrorNoSuchKey = errors.New("no such key")

// Put data in Stash
func (s *Stash) Put(key string, doc interface{}) error {
	s.Lock()
	s.m[key] = doc
	s.Unlock()

	// TODO: transaction log
	return nil
}

// Get data from Stash
func (s *Stash) Get(key string) (interface{}, error) {
	s.RLock()
	doc, ok := s.m[key]
	s.RUnlock()
	if !ok {
		return nil, ErrorNoSuchKey
	}

	return doc, nil
}

// Delete data from stash
func (s *Stash) Delete(key string) error {
	s.Lock()
	delete(s.m, key)
	s.Unlock()

	// TODO: transaction log
	return nil
}

// Restore not implement
// TODO: restore from transaction log
func (s *Stash) Restore() error {
	return nil
}

/*
func (s *Stash) Restore() error {
	var err error

	events, errors := store.transact.ReadEvents()
	count, ok, e := 0, true, Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = store.Delete(e.Key)
				count++
			case EventPut: // Got a PUT event!
				err = store.Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	store.transact.Run()

	go func() {
		for err := range store.transact.Err() {
			log.Print(err)
		}
	}()

	return err
}*/
