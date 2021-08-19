package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTasks() []*Task {
	return []*Task{
		{Id: 0, Title: "Task One", Description: "Important task"},
		{Id: 1, Title: "Task Two", Description: "Another important task"},
	}
}

func assertTaskEqual(t *testing.T, expectedTask *Task, actualTask *Task) {
	assertEqual(t, expectedTask.Id, actualTask.Id)
	assertEqual(t, expectedTask.Title, actualTask.Title)
	assertEqual(t, expectedTask.Description, actualTask.Description)
}

func assertTasksEqual(t *testing.T, expectedTasks []*Task, actualTasks []*Task) {
	assertEqual(t, len(expectedTasks), len(actualTasks))
	for i := range expectedTasks {
		assertTaskEqual(t, expectedTasks[i], actualTasks[i])
	}
}

func TestTasks(t *testing.T) {
	t.Run("get all tasks", func(t *testing.T) {
		expectedTasks := createTasks()
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(expectedTasks))
		request, _ := http.NewRequest(http.MethodGet, "/tasks/", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		actualTasks := []*Task{}
		assertNoError(t, json.Unmarshal(response.Body.Bytes(), &actualTasks))
		assertTasksEqual(t, expectedTasks, actualTasks)
	})

	t.Run("get task by id", func(t *testing.T) {
		expectedTasks := createTasks()
		expectedTask := expectedTasks[0]
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(expectedTasks))
		request, _ := http.NewRequest(http.MethodGet, "/tasks/0", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		actualTask := Task{}
		assertNoError(t, json.Unmarshal(response.Body.Bytes(), &actualTask))
		assertTaskEqual(t, expectedTask, &actualTask)
	})
}
