package main

import (
	"github.com/freahs/lunch-server"
	"github.com/freahs/lunch-server/loaders"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func loadStore(filename string) *lunch_server.Store {
	dir := func() string {
		if d, err := os.UserConfigDir(); err == nil {
			d = filepath.Join(d, "lunch-server")
			if err = os.MkdirAll(d, 0755); err == nil {
				return d
			}
		}
		return os.TempDir()
	}()
	path := filepath.Join(dir, filename)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return lunch_server.NewStore()
	}
	if info.IsDir() {
		log.Fatalf("%v is a directory", path)
	}
	store, err := lunch_server.LoadStore(path)
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := lunch_server.NewStore()
	prime := loaders.NewPrime()
	err := prime.Load(store)
	if err != nil {
		log.Fatal(err)
	}
	Router := mux.NewRouter().StrictSlash(true)
	APIRouter := Router.PathPrefix("/api/v1").Subrouter()
	api := NewAPIServer(store, APIRouter)
	log.Printf("starting api server @ port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, api))
}