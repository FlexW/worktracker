package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TaskWithoutId struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TimeIntervals []TimeInterval `json:"time_intervals"`
}

func taskWithoutIdToTask(taskWithoutId *TaskWithoutId, id int) *Task {
	return &Task{
		Id: id,
		Title: taskWithoutId.Title,
		Description: taskWithoutId.Description,
		TimeIntervals: taskWithoutId.TimeIntervals,
	}
}

func createTasks() []*Task {
	endTime := time.Now().Add(time.Hour)

	return []*Task{
		{
			Id: 0,
			Title: "Task One",
			Description: "Important task",
			TimeIntervals: []TimeInterval{{time.Now(), &endTime}},
		},
		{
			Id: 1,
			Title: "Task Two",
			Description: "Another important task",
			TimeIntervals: []TimeInterval{{time.Now(), nil}},
		},
	}
}

func assertTasksInactive(t *testing.T, tasks []*Task) {
	for _, task := range tasks {
		for _, timeInterval := range task.TimeIntervals {
			assert.NotNil(t, timeInterval.EndTime)
		}
	}
}

func assertTimeIntervalEqual(t *testing.T, expected TimeInterval, actual TimeInterval) {
	assert.True(t, expected.StartTime.Equal(*actual.EndTime))
	if expected.EndTime != nil {
		assert.True(t, actual.EndTime != nil)
		assert.True(t, expected.EndTime.Equal(*actual.EndTime))
	} else {
		assert.True(t, actual.EndTime == nil)
	}
}

func assertTimeIntervalsEqual(t *testing.T, expected []TimeInterval, actual []TimeInterval) {
	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assertTimeIntervalEqual(t, expected[i], actual[i])
	}
}

func assertTaskEqual(t *testing.T, expectedTask *Task, actualTask *Task) {
	assert.Equal(t, expectedTask.Id, actualTask.Id)
	assert.Equal(t, expectedTask.Title, actualTask.Title)
	assert.Equal(t, expectedTask.Description, actualTask.Description)
}

func assertTasksEqual(t *testing.T, expectedTasks []*Task, actualTasks []*Task) {
	assert.Equal(t, len(expectedTasks), len(actualTasks))
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
		json.Unmarshal(response.Body.Bytes(), &actualTasks)
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
		json.Unmarshal(response.Body.Bytes(), &actualTask)
		assertTaskEqual(t, expectedTask, &actualTask)
	})

	t.Run("create new task", func(t *testing.T) {
		newTask := TaskWithoutId{
			Title: "Some Task",
			Description: "Description",
			TimeIntervals: []TimeInterval{{time.Now(), nil}},
		}
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(make([]*Task,0)))
		newTaskJson, _ := json.Marshal(&newTask)
		request, _ := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewBuffer(newTaskJson))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskInStore := worktrackerServer.store.GetAllTasks()[0]
		assertTaskEqual(t, taskWithoutIdToTask(&newTask, taskInStore.Id), taskInStore)
	})

	t.Run("create new task without end time then finish current active task", func(t *testing.T) {
		newTask := Task{
			Id: 2,
			Title: "Task One",
			Description: "Important task",
			TimeIntervals: []TimeInterval{{time.Now(), nil}},
		}
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(createTasks()))
		newTaskJson, _ := json.Marshal(&newTask)
		request, _ := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewBuffer(newTaskJson))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		tasksInStore := worktrackerServer.store.GetAllTasks()
		for _, taskInStore := range tasksInStore {
			if taskInStore.Id != newTask.Id {
				for _, timestamps := range taskInStore.TimeIntervals {
					assert.NotNil(t, timestamps.EndTime)
				}
			} else {
				assert.Nil(t, taskInStore.TimeIntervals[0].EndTime)
			}
		}
	})

	t.Run("stop current task", func(t *testing.T) {
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(createTasks()))
		request, _ := http.NewRequest(http.MethodPost, "/tasks/stop", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		tasksInStore := worktrackerServer.store.GetAllTasks()
		assertTasksInactive(t, tasksInStore)
	})

	t.Run("start stopped task again", func(t *testing.T) {
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(createTasks()))
		timeNow := time.Now()
		requestData := []byte(fmt.Sprintf(`{
            "start_time": "%v"
        }`, timeNow.Format(time.RFC3339)))
		request, _ := http.NewRequest(http.MethodPost, "/tasks/0", bytes.NewBuffer(requestData))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskInStore := worktrackerServer.store.GetTaskById(0)
		assert.NotNil(t, taskInStore.TimeIntervals[0].EndTime)
		assert.Equal(t, timeNow.Truncate(time.Second), taskInStore.TimeIntervals[1].StartTime)
	})
}
