package stash

import (
	"encoding/json"
	"sync"
)

// TODO: encrypt data o

import (
	"errors"
)

// Stash in-memory storage
type Stash struct {
	sync.RWMutex
	m map[string]string
}

// NewStash constructor
func NewStash() *Stash {
	return &Stash{
		m: make(map[string]string),
	}
}

// Read implement io.Reader
func (s *Stash) Read(p []byte) (n int, err error) {
	err = json.Unmarshal(p, &s.m)
	return len(p), err
}

// Write implement io.Writer
func (s *Stash) Write(p []byte) (n int, err error) {
	p, err = json.Marshal(s.m)
	return len(p), err
}

// ErrorNoSuchKey predefined error
var ErrorNoSuchKey = errors.New("no such key")

// Put data in Stash
func (s *Stash) Put(key, data string) error {
	s.Lock()
	s.m[key] = data
	s.Unlock()

	// TODO: transaction log
	return nil
}

// Get data from Stash
func (s *Stash) Get(key string) (string, error) {
	s.RLock()
	value, ok := s.m[key]
	s.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
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
