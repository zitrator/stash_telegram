package main

import (
	_ "errors"
	"fmt"
	"sync"
)

type Stash struct {
	sync.RWMutex
	m map[string]string
}

// Stash entry point
func main() {
	fmt.Println("vim-go")
}
