package main

import (
	"net/http"
	"time"

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
	w.router.POST("/tasks/stop", w.handleStopTasks)
	w.router.GET("/tasks/:id", w.handleTaskById)
	w.router.POST("/tasks/:id", w.handleStartTask)
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

func (w *WorktrackerServer) setTaskInactive(task *Task) {
	for i := range task.TimeIntervals {
		if task.TimeIntervals[i].EndTime == nil {
			timeNow := time.Now()
			task.TimeIntervals[i].EndTime = &timeNow
			// TODO: Make sure task gets saved to disk. Capture this with e2e test
		}
	}
}

type startTask struct {
	StartTime string `json:"start_time"`
}

func (w *WorktrackerServer) handleStartTask(c *gin.Context) {
	taskId := c.GetInt("id")
	var startTimeStr startTask
	if err := c.BindJSON(&startTimeStr); err != nil {
		return
	}
	startTime, err := time.Parse(time.RFC3339, startTimeStr.StartTime)
	if err != nil {
		c.Err()
	}
	w.setAllTasksInactive()
	task := w.store.GetTaskById(taskId)
	task.TimeIntervals = append(task.TimeIntervals, TimeInterval{StartTime: startTime, EndTime: nil})
	w.store.UpdateTask(task)
	c.IndentedJSON(http.StatusOK, task)
}

func (w *WorktrackerServer) handleStopTasks(c *gin.Context) {
	w.setAllTasksInactive()
}

func (w *WorktrackerServer) setAllTasksInactive() {
	tasks := w.store.GetAllTasks()
	for _, task := range tasks {
		w.setTaskInactive(task)
	}
}

func (w *WorktrackerServer) handleNewTask(c *gin.Context) {
	var newTask Task

	if err:= c.BindJSON(&newTask); err != nil {
		return
	}
	w.setAllTasksInactive()
	newTask.Id = w.store.InsertTask(&newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}

func (w *WorktrackerServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	w.router.ServeHTTP(responseWriter, request)
}
