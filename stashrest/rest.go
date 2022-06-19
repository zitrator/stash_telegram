package stashrest

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitrator/stash_telegram/stash"
)

type stashRest struct {
	stash *stash.Stash
}

func NewStashRest(st *stash.Stash) *stashRest {
	return &stashRest{stash: st}
}

// Start initialize and start the mux.Router
func (sr *stashRest) Start() error {
	router := mux.NewRouter()
	router.Use(sr.loggingMiddleware)

	router.HandleFunc("/s/{key}", sr.keyValueGetHandler).Methods("GET")
	router.HandleFunc("/s/{key}", sr.keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/s/{key}", sr.keyValueDeleteHandler).Methods("DELETE")

	// todo: add support for multi stashes
	//router.HandleFunc("/s/{stash}/{key}", sr.keyValueGetHandler).Methods("GET")
	//router.HandleFunc("/s/{stash}/{key}", sr.keyValuePutHandler).Methods("PUT")
	//router.HandleFunc("/s/{stash}/{key}", sr.keyValueDeleteHandler).Methods("DELETE")

	//router.HandleFunc("/s/{stash}/{key}", sr.keyValueGetHandler).Methods("GET")
	//router.HandleFunc("/s/{stash}/{key}", sr.keyValuePutHandler).Methods("PUT")
	//router.HandleFunc("/s/{stash}/{key}", sr.keyValueDeleteHandler).Methods("DELETE")

	router.HandleFunc("/s", sr.notAllowedHandler)
	router.HandleFunc("/s/{key}", sr.notAllowedHandler)

	// todo: port from environment variable
	return http.ListenAndServe(":8080", router)
}

// loggingMiddleware the handler for logging requests
func (sr *stashRest) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// notAllowedHandler type method stashRest
func (sr *stashRest) notAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

// keyValuePutHandler type method stashRest
func (sr *stashRest) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := ioutil.ReadAll(r.Body)
	defer wrapErrors(r.Body.Close())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sr.stash.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Printf("PUT key=%s value=%s\n", key, string(value))
}

// keyValueGetHandler type method stashRest
func (sr *stashRest) keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := sr.stash.Get(key)
	if errors.Is(err, stash.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(value))
	wrapErrors(err)

	log.Printf("GET key=%s\n", key)
}

// keyValueDeleteHandler type method stashRest
func (sr *stashRest) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := sr.stash.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DELETE key=%s\n", key)
}

func wrapErrors(err error) {
	if err != nil {
		log.Println(err)
	}
}
