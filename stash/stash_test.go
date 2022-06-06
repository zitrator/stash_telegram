package stash

import (
	"testing"
	"testing/quick"
)

var config quick.Config

// TODO: one stash on all test

func TestStash_Put(t *testing.T) {
	stash := NewStash()
	// store data
	f := func(key, val string) bool {
		ok := stash.Put(key, val)
		if ok != nil {
			return false
		}
		return true
	}
	if err := quick.Check(f, &config); err != nil {
		t.Error(err)
	}
}

func TestStash_Get(t *testing.T) {
	stash := NewStash()
	// store data
	put := func(key, val string) bool {
		ok := stash.Put(key, val)
		if ok != nil {
			return false
		}
		return true
	}
	// get data
	get := func(key, val string) bool {
		res, ok := stash.Get(key)
		if ok != nil || res != val {
			return false
		}
		return true
	}

	if err := quick.CheckEqual(put, get, &config); err != nil {
		t.Error(err)
	}
}

func TestStash_Delete(t *testing.T) {
	stash := NewStash()
	// store data
	put := func(key, val string) bool {
		ok := stash.Put(key, val)
		if ok != nil {
			return false
		}
		return true
	}
	// del data
	del := func(key, val string) bool {
		ok := stash.Delete(key)
		if ok != nil {
			return false
		}
		_, ok = stash.Get(key)
		if ok == nil {
			return false
		}
		return true
	}

	if err := quick.CheckEqual(put, del, &config); err != nil {
		t.Error(err)
	}
}
