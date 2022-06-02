package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Stash in-memory storage
// todo: transaction log
// todo: encrypt data on disk
type Stash struct {
	sync.RWMutex
	m map[string]string
}

var jsonData = []byte("{ \"name\":\"patrik\", \"age\":10 }")

func main() {
	var f interface{}
	err := json.Unmarshal(jsonData, &f)

	if err != nil {
		fmt.Println("err")
	}

	mj := f.(map[string]interface{})
	fmt.Println(f)
	fmt.Println(mj["name"])
	fmt.Println(mj["age"])
	fmt.Println(mj["none"])
}
