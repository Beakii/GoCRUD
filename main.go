package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStoreage()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":80", store)
	server.Run()
}
