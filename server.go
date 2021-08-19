package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type WorktrackerServer struct {
	store WorktrackerStore
	router *gin.Engine
}

func NewWorktrackerServer(store WorktrackerStore) *WorktrackerServer {
	w := &WorktrackerServer{store, gin.Default()}
	w.router.GET("/tasks/", w.handleTasks)
	w.router.POST("/tasks/", w.handleNewTask)
	w.router.GET("/tasks/:id", w.handleTaskById)
	return w;
}

func (w *WorktrackerServer) handleTasks(c *gin.Context) {
	tasks := w.store.GetAllTasks()
	c.IndentedJSON(http.StatusOK, tasks)
}

func (w *WorktrackerServer) handleTaskById(c *gin.Context) {
	taskId := c.GetInt("id")
	task := w.store.GetTaskById(taskId)
	c.IndentedJSON(http.StatusOK, task)
}

func (w *WorktrackerServer) handleNewTask(c *gin.Context) {
	var newTask Task

	if err:= c.BindJSON(&newTask); err != nil {
		return
	}
	w.store.InsertTask(&newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func (w *WorktrackerServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	w.router.ServeHTTP(responseWriter, request)
}
