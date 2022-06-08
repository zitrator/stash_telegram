package stash

import (
	"fmt"
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

func TestStash_Read(t *testing.T) {
	s := []byte("{\"key1\":\"val1\",\"key2\":\"100\"}")
	_, err := stash.Read(s)
	if err != nil {
		t.Error(err)
	}
	if stash.m["key1"] != "val1" && stash.m["key2"] != "100" {
		t.Error("Read wrong data")
	}
}

func TestStash_Write(t *testing.T) {
	err := stash.Put("key1", "value1")
	if err != nil {
		t.Error(err)
	}
	err = stash.Put("key2", "value2")
	if err != nil {
		t.Error(err)
	}
	_, err = fmt.Fprintln(stash)
	if err != nil {
		t.Error(err)
	}
}
