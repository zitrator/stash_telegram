package main

import (
	"github.com/zitrator/stash_telegram/stash"
	"github.com/zitrator/stash_telegram/stashrest"
	"log"
)

type Front interface {
	Start(s *stash.Stash) error
}

func main() {
	log.Println("Starting stashes ...")
	log.Fatal(stashrest.NewStashRest().Start("my_first_db"))
}
