package stash

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
	return &Stash{m: make(map[string]string)}
}

// marshal read from Stash into p
func (s *Stash) marshal() (p []byte, err error) {
	p, err = json.Marshal(s.m)
	return p, err
}

// unmarshal from p to the Stash
func (s *Stash) unmarshal(p []byte) (err error) {
	err = json.Unmarshal(p, &s.m)
	return err
}

// Backup data
func (s *Stash) Backup(w io.Writer) error {
	p, err := s.marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(p)
	return nil
}

func (s *Stash) Restore(r io.Reader) (err error) {
	p, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = s.unmarshal(p)
	return err
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

// RestoreFromTranLog
// TODO: restore from transaction log
func (s *Stash) RestoreFromTranLog() error {
	return errors.New("not implemented")
}
