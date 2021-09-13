package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	server := NewWorktrackerServer(
		NewInMemoryWorktrackerStore(make([]*Task, 0), map[int][]*TimeInterval{}))

	log.Err(http.ListenAndServe(":12345", server))
}
