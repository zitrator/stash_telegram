package stash

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// todo: encrypt data
// todo: gzip data

const (
	systemStashName      = ".system"
	stashesFolderVarName = "STASHES_ROOT"
	syncDataEvery        = 10 * time.Second
)

type Database struct {
	sync.RWMutex
	once    sync.Once
	stashes map[string]*Stash
	folder  string
	ticker  *time.Ticker
	quit    chan interface{}
}

// database one instance for all stashes
var database *Database

// init database instance and restore stashes
func init() {
	folderName := os.Getenv(stashesFolderVarName)
	if folderName == "" {
		log.Fatal(stashesFolderVarName + " variable not set")
	}

	var err error
	database, err = initNewDatabase(folderName)
	if err != nil {
		log.Fatal(err)
	}
	database.ticker = time.NewTicker(syncDataEvery)
	go tickerFunc()
}

// tickerFunc save data to disk every syncDataEvery
func tickerFunc() {
	for {
		select {
		case <-database.ticker.C:
			database.RLock()
			cm := make(map[string]*Stash)
			for key, ptr := range database.stashes {
				cm[key] = ptr
			}
			database.RUnlock()
			for fn, s := range cm {
				saveStash(s, database.folder+"/"+fn)
			}
		case <-database.quit:
			database.ticker.Stop()
			return
		}
	}
}

// saveStash save the stash (*st) to file (fn)
func saveStash(st *Stash, fn string) {
	if !st.ch {
		return
	}
	fo, err := os.Create(fn)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			log.Println(err)
		}
	}()

	err = st.Backup(fo)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Stash %st was saved", fn)
	}
}

// GetDatabase return the initialized instance
func GetDatabase() *Database {
	return database
}

// GetStash pointer, the stash will be created if it doesn't exist
func (db *Database) GetStash(id string) *Stash {
	// todo: block the system stash
	db.once.Do(func() {
		if stash := db.stashes[id]; stash == nil {
			db.stashes[id] = NewStash()
		}
	})
	return db.stashes[id]
}

// initNewDatabase initialize and return the new database
func initNewDatabase(folderName string) (*Database, error) {
	if database != nil {
		return database, errors.New("the database has already been initialized, this call has been ignored")
	}

	if _, err := os.Stat(folderName); err != nil {
		if os.IsNotExist(err) {
			// the folder doesn't exist
			if err = os.MkdirAll(folderName, 0750); err != nil {
				log.Fatal(err)
			}
		} else {
			// the unknown fatal error
			log.Fatal(err)
		}
	}

	database := Database{folder: folderName, stashes: make(map[string]*Stash)}
	if files, err := ioutil.ReadDir(folderName); err == nil {
		for _, file := range files {
			st := NewStash()
			if data, err := os.ReadFile(folderName + "/" + file.Name()); err == nil {
				if err = st.Restore(bytes.NewReader(data)); err != nil {
					log.Println(err)
				}
				database.stashes[file.Name()] = st
				log.Println("the file ", file.Name(), " was restored")
			} else {
				log.Println("the file ", file.Name(), " was ignored: ", err)
			}
		}
	} else {
		log.Fatal(err)
	}
	return &database, nil
}
