package stash

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"
)

// todo: multiple stash support
// todo: encrypt data
// todo: transaction log

import (
	"errors"
)

// Stash in-memory key value storage
type Stash struct {
	sync.RWMutex
	m       map[string]string
	id      string
	changed bool
}

// NewStash constructor
func NewStash(id string) *Stash {
	return &Stash{m: make(map[string]string), id: id, changed: false}
}

// marshal read from Stash into p
func (s *Stash) marshal() (p []byte, err error) {
	p, err = json.Marshal(s.m)
	return p, err
}

// unmarshal from p to the Stash
func (s *Stash) unmarshal(p []byte) error {
	err := json.Unmarshal(p, &s.m)
	return err
}

// Backup data
func (s *Stash) Backup(w io.Writer) error {
	s.Lock()
	p, err := s.marshal()
	if err != nil {
		s.Unlock()
		return err
	}
	s.changed = false
	s.Unlock()
	_, err = w.Write(p)

	return nil
}

func (s *Stash) Restore(r io.Reader) error {
	p, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.Lock()
	err = s.unmarshal(p)
	s.changed = false
	s.Unlock()

	return err
}

// ErrorNoSuchKey predefined error
var ErrorNoSuchKey = errors.New("no such key")

// Put data in Stash
func (s *Stash) Put(key, data string) error {
	s.Lock()
	s.m[key] = data
	s.Unlock()
	s.changed = true

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
	s.changed = true

	return nil
}
