package main

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	Title       string
	Description string
}

type WorktrackerStore interface{
	GetAllTasks() []*Task
}

type WorktrackerServer struct {
	store WorktrackerStore
}

func (w *WorktrackerServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	tasks := w.store.GetAllTasks()
	tasksJson, _ := json.Marshal(&tasks)
	responseWriter.Write(tasksJson)
}
