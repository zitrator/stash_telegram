package stash

import (
	"testing"
	"testing/quick"
)

var config quick.Config
var stash *Stash

func init() {
	stash = NewStash()
}

func TestStash_Put(t *testing.T) {
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
	// store the data
	put := func(key, val string) bool {
		ok := stash.Put(key, val)
		if ok != nil {
			return false
		}
		return true
	}
	// get the data
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
	// store the data
	put := func(key, val string) bool {
		ok := stash.Put(key, val)
		if ok != nil {
			return false
		}
		return true
	}
	// delete the data
	del := func(key, val string) bool {
		ok := stash.Delete(key)
		if ok != nil {
			return false
		}
		// check that the data is deleted
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
