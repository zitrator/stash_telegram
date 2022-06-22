package stashrest

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitrator/stash_telegram/stash"
)

type stashRest struct {
	database *stash.Database
}

func NewStashRest(db *stash.Database) *stashRest {
	if db == nil {
		log.Fatal("nil database pointer")
	}
	return &stashRest{database: db}
}

// Start initialize and start the mux.Router
func (sr *stashRest) Start() error {
	router := mux.NewRouter()
	router.Use(sr.logging)

	router.HandleFunc("/s/{stash}", sr.stashGetHandler).Methods("GET")
	router.HandleFunc("/s/{stash}", sr.stashPutHandler).Methods("PUT")
	router.HandleFunc("/s/{stash}", sr.stashDeleteHandler).Methods("DELETE")

	router.HandleFunc("/s/{stash}/{key}", sr.dataGetHandler).Methods("GET")
	router.HandleFunc("/s/{stash}/{key}", sr.dataPutHandler).Methods("PUT")
	router.HandleFunc("/s/{stash}/{key}", sr.dataDeleteHandler).Methods("DELETE")

	router.HandleFunc("/s", sr.notAllowedHandler)
	router.HandleFunc("/s/{stash}", sr.notAllowedHandler)
	router.HandleFunc("/s/{stash}/{key}", sr.notAllowedHandler)

	// todo: port from environment variable
	return http.ListenAndServe(":8080", router)
}

// logging the handler for logging requests
func (sr *stashRest) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// notAllowedHandler the handler of incorrect requests
func (sr *stashRest) notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

// dataPutHandler data putting request handler
func (sr *stashRest) dataPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if st, ok := sr.database.Get(vars["stash"]); ok {
		// the stash ok
		value, err := ioutil.ReadAll(r.Body)
		defer wrapErr(r.Body.Close())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = st.Put(vars["key"], string(value))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		log.Printf("PUT vars=%v value=%s\n", vars, string(value))
	} else {
		// the stash not found
		http.Error(w, "no such stash", http.StatusNotFound)
	}
}

// dataGetHandler data getting request handler
func (sr *stashRest) dataGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if st, ok := sr.database.Get(vars["stash"]); ok {
		// the stash is fine
		if value, err := st.Get(vars["key"]); err == nil {
			// stash.Get is fine
			if _, err = w.Write([]byte(value)); err == nil {
				log.Printf("GET vars=%v value=%s\n", vars, value)
			} else {
				wrapErr(err)
			}
		} else {
			if errors.Is(err, stash.ErrorNoSuchKey) {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		// the stash not found
		http.Error(w, "no such stash", http.StatusNotFound)
	}
}

// dataDeleteHandler data deletion request handler
func (sr *stashRest) dataDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if st, ok := sr.database.Get(vars["stash"]); ok {
		// the stash is fine
		if err := st.Delete(vars["key"]); err == nil {
			log.Printf("DELETE stash=%v\n", vars)
		} else {
			if errors.Is(err, stash.ErrorNoSuchKey) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		// the stash not found
		http.Error(w, "no such stash", http.StatusNotFound)
	}
}

// stashGetHandler stash getting request handler
func (sr *stashRest) stashGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if st, ok := sr.database.Get(vars["stash"]); ok {
		// the stash is fine
		var buf bytes.Buffer
		if err := st.Backup(&buf); err == nil {
			if _, err := w.Write(buf.Bytes()); err == nil {
				log.Printf("GET vars=%v %s\n", vars, buf.String())
			} else {
				wrapErr(err)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// the stash not found
		http.Error(w, stash.ErrorNoSuchStash.Error(), http.StatusNotFound)
	}
}

// stashPutHandler stash creating request handler
func (sr *stashRest) stashPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, ok := sr.database.Get(vars["stash"]); ok {
		log.Printf("GET vars=%v\n", vars)
	}
}

// stashDeleteHandler stash deleting request handler
func (sr *stashRest) stashDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := sr.database.Delete(vars["stash"]); err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE vars=%v\n", vars)
}

func wrapErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
