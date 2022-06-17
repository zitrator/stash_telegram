package stash

import (
	"bytes"
	"log"
	"os"
	"sync"
	"time"
)

const systemStashName = ".system"
const stashesFolderVarName = "STASHES_ROOT"
const syncDataEvery = 10 * time.Second

// todo: singleton
// todo: sync data in goroutine

type Database struct {
	sync.RWMutex
	once    sync.Once
	stashes map[string]*Stash
	folder  string
	ticker  *time.Ticker
	quit    chan interface{}
}

// instance one instance for all stashes
var instance *Database

func init() {
	folderName := os.Getenv(stashesFolderVarName)
	if folderName == "" {
		log.Fatal(stashesFolderVarName + " variable not set")
	}

	var err error
	instance, err = initNewDatabase(folderName)
	if err != nil {
		log.Fatal(err)
	}
	instance.ticker = time.NewTicker(syncDataEvery)
	go tickerFunc()
}

func tickerFunc() {
	for {
		select {
		case <-instance.ticker.C:
			instance.RLock()
			cs := make(map[string]*Stash)
			for key, ptr := range instance.stashes {
				cs[key] = ptr
			}
			instance.RUnlock()
			for key, ptr := range cs {
				saveStash(instance.folder+"/"+key, ptr)
			}
		case <-instance.quit:
			instance.ticker.Stop()
			return
		}
	}
}

func saveStash(key string, ptr *Stash) {
	if !ptr.changed {
		return
	}
	fo, err := os.Create(key)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			log.Println(err)
		}
	}()

	err = ptr.Backup(fo)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Stash %s was saved", key)
	}
}

// GetDatabase init new database
func GetDatabase() *Database {
	return instance
}

// GetStash pointer, the stash will be created if it doesn't exist
func (db *Database) GetStash(id string) *Stash {
	// todo: block the system stash
	db.once.Do(func() {
		if stash := db.stashes[id]; stash == nil {
			db.stashes[id] = NewStash(id)
		}
	})
	return db.stashes[id]
}

// initNewDatabase
func initNewDatabase(folderName string) (*Database, error) {

	if _, err := os.Stat(folderName); err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(folderName, 0750); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	database := Database{folder: folderName, stashes: make(map[string]*Stash)}
	stashPtr := NewStash(systemStashName)
	if data, err := os.ReadFile(systemStashName); err == nil {
		if err = stashPtr.Restore(bytes.NewReader(data)); err != nil {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}
	database.stashes[systemStashName] = stashPtr
	return &database, nil
}
