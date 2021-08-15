package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTasks() []*Task {
	return []*Task{
		{Title: "Task One", Description: "Important task"},
		{Title: "Task Two", Description: "Another important task"},
	}
}

func assertTaskEqual(t *testing.T, expectedTask *Task, actualTask *Task) {
	assertEqual(t, expectedTask.Title, actualTask.Title)
	assertEqual(t, expectedTask.Description, actualTask.Description)
}

func assertTasksEqual(t *testing.T, expectedTasks []*Task, actualTasks []*Task) {
	assertEqual(t, len(expectedTasks), len(actualTasks))
	for i := range expectedTasks {
		assertTaskEqual(t, expectedTasks[i], actualTasks[i])
	}
}

type FakeWorktrackerStore struct {
	tasks []*Task
}

func (s *FakeWorktrackerStore) GetAllTasks() []*Task {
	return s.tasks
}

func TestGetTasks(t *testing.T) {
	t.Run("get all tasks", func(t *testing.T) {
		expectedTasks := createTasks()
		worktrackerServer := &WorktrackerServer{store: &FakeWorktrackerStore{expectedTasks}}
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		actualTasks := []*Task{}
		assertNoError(t, json.Unmarshal(response.Body.Bytes(), &actualTasks))
		assertTasksEqual(t, expectedTasks, actualTasks)
	})
}
