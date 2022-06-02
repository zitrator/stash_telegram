package main

import (
	"encoding/json"
	"fmt"
	"github.com/zitrator/stash_telegram/stash"
)

var jsonData = []byte("{ \"name\":\"patrik\", \"age\":10 }")

func main() {
	myStash := stash.Stash{}
	fmt.Println(myStash)

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
