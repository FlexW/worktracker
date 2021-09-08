package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewWorktrackerServer(NewInMemoryWorktrackerStore(make([]*Task, 0), map[int][]*TimeInterval{}))
	log.Fatal(http.ListenAndServe(":12345", server))
}
