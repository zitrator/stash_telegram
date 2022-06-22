package stash

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
)

// todo: transaction log

import (
	"errors"
)

var ErrorNoSuchKey = errors.New("no such key")

// Stash in-memory key value storage
type Stash struct {
	sync.RWMutex
	// a map for the stash
	m map[string]string
	// true if the data has been changed
	ch bool
}

// NewStash constructor
func NewStash() *Stash {
	return &Stash{m: make(map[string]string), ch: false}
}

// marshal Stash => p
func (s *Stash) marshal() (p []byte, err error) {
	p, err = json.Marshal(s.m)
	return p, err
}

// unmarshal p => Stash
func (s *Stash) unmarshal(p []byte) error {
	err := json.Unmarshal(p, &s.m)
	return err
}

// Backup the stash data
func (s *Stash) Backup(w io.Writer) error {
	s.Lock()
	p, err := s.marshal()
	if err != nil {
		s.Unlock()
		return err
	}
	s.ch = false
	s.Unlock()

	_, err = w.Write(p)
	return err
}

// Restore the stash data
func (s *Stash) Restore(r io.Reader) error {
	p, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.Lock()
	err = s.unmarshal(p)
	s.ch = false
	s.Unlock()

	return err
}

// Put the data in the Stash
func (s *Stash) Put(key, data string) error {
	s.Lock()
	s.m[key] = data
	s.Unlock()
	s.ch = true

	return nil
}

// Get the data from the Stash
func (s *Stash) Get(key string) (string, error) {
	s.RLock()
	value, ok := s.m[key]
	s.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

// Delete the data from the stash
func (s *Stash) Delete(key string) error {
	s.Lock()
	delete(s.m, key)
	s.Unlock()
	s.ch = true

	return nil
}
