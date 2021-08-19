package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewWorktrackerServer(NewInMemoryWorktrackerStore(make([]*Task,0)))
	log.Fatal(http.ListenAndServe(":12345", server))
}
