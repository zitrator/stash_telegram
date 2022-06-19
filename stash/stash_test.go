package stash

import (
	"os"
	"reflect"
	"testing"
	"testing/quick"
)

const testDatabase = "test_db"

type testContext struct {
	fileName string
	test01m  map[string]string
	test01b  []byte
}

var config quick.Config
var context testContext

func init() {
	context.fileName = "/tmp/stash.json"
	context.test01m = map[string]string{"key1": "value1", "key2": "100"}
	context.test01b = []byte("{\"key1\":\"value1\",\"key2\":\"100\"}")
}

func TestStash_Put(t *testing.T) {
	stash := NewStash()
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
	stash := NewStash()
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

func TestStash_marshal(t *testing.T) {
	stash := NewStash()
	var err error
	for key, val := range context.test01m {
		err = stash.Put(key, val)
		if err != nil {
			t.Error(err)
		}
	}
	var p []byte
	p, err = stash.marshal()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(p, context.test01b) {
		t.Error(p, " != ", context.test01b)
	}
}

func TestStash_unmarshal(t *testing.T) {
	stash := NewStash()
	err := stash.unmarshal(context.test01b)
	if err != nil {
		t.Error(err)
	}
	for key, val := range context.test01m {
		if stash.m[key] != val {
			t.Error(key, stash.m[key], val)
		}
	}
}

func TestStash_Backup(t *testing.T) {
	f, err := os.OpenFile(context.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)

	stash := NewStash()
	for key, val := range context.test01m {
		err = stash.Put(key, val)
		if err != nil {
			t.Error(err)
		}
	}
	err = stash.Backup(f)
	if err != nil {
		t.Error(err)
	}
}

func TestStash_Restore(t *testing.T) {
	f, err := os.OpenFile(context.fileName, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Error(err)
		}
	}(f)

	TestStash_Backup(t)
	stash := NewStash()
	err = stash.Restore(f)
	if err != nil {
		t.Error(err)
	}
	for key, val := range context.test01m {
		if stash.m[key] != val {
			t.Error(key, stash.m[key], val)
		}
	}
}
