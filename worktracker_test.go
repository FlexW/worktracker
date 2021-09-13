package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

type NewTaskWithoutId struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"startTime"`
	EndTime     *time.Time `json:"endTime"`
}

func newTaskWithoutIdToTask(newTaskWithoutId *NewTaskWithoutId, id int) *Task {
	var duration time.Duration
	duration = 0
	active := true
	if newTaskWithoutId.EndTime != nil {
		duration = newTaskWithoutId.EndTime.Sub(*newTaskWithoutId.EndTime)
		active = false
	}
	return &Task{
		Id:          id,
		Title:       newTaskWithoutId.Title,
		Description: newTaskWithoutId.Description,
		Duration:    duration,
		Active:      active,
	}
}

func createTasks() []*Task {
	return []*Task{
		{
			Id:          0,
			Title:       "Task One",
			Description: "Important task",
			Duration:    1994,
			Active:      false,
		},
		{
			Id:          1,
			Title:       "Task Two",
			Description: "Another important task",
			Duration:    207,
			Active:      false,
		},
	}
}

func assertTasksInactive(t *testing.T, tasks []*Task) {
	for _, task := range tasks {
		assert.False(t, task.Active)
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
	assert.Equal(t, expectedTask.Active, actualTask.Active)
	assert.Equal(t, expectedTask.Duration, actualTask.Duration)
}

func assertTasksEqual(t *testing.T, expectedTasks []*Task, actualTasks []*Task) {
	assert.Equal(t, len(expectedTasks), len(actualTasks))
	for i := range expectedTasks {
		found := false
		for j := range actualTasks {
			if expectedTasks[i].Id == actualTasks[j].Id {
				assertTaskEqual(t, expectedTasks[i], actualTasks[j])
				found = true
			}
		}
		assert.True(t, found)
	}
}

func TestTasks(t *testing.T) {
	t.Run("get all tasks", func(t *testing.T) {
		expectedTasks := createTasks()
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(expectedTasks))
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
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(expectedTasks))
		request, _ := http.NewRequest(http.MethodGet, "/tasks/0", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		actualTask := Task{}
		json.Unmarshal(response.Body.Bytes(), &actualTask)
		assertTaskEqual(t, expectedTask, &actualTask)
	})

	t.Run("create new active task", func(t *testing.T) {
		newTask := NewTaskWithoutId{
			Title:       "Some Task",
			Description: "Description",
			StartTime:   time.Now(),
			EndTime:     nil,
		}
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(make([]*Task, 0)))
		newTaskJson, _ := json.Marshal(&newTask)
		request, _ := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewBuffer(newTaskJson))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskInStore := worktrackerServer.store.GetAllTasks()[0]
		assertTaskEqual(t, newTaskWithoutIdToTask(&newTask, taskInStore.Id), taskInStore)
	})

	t.Run("create new inactive task", func(t *testing.T) {
		startTime := time.Now()
		endTime := startTime.Add(time.Hour)
		newTask := NewTaskWithoutId{
			Title:       "Some Task",
			Description: "Description",
			StartTime:   startTime,
			EndTime:     &endTime,
		}
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(make([]*Task, 0)))
		newTaskJson, _ := json.Marshal(&newTask)
		request, _ := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewBuffer(newTaskJson))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskInStore := worktrackerServer.store.GetAllTasks()[0]
		assertTaskEqual(t, newTaskWithoutIdToTask(&newTask, taskInStore.Id), taskInStore)
	})

	t.Run("create new active task then finish current active task", func(t *testing.T) {
		newTask := Task{
			Id:          2,
			Title:       "Task One",
			Description: "Important task",
			Duration:    2008,
			Active:      true,
		}
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(createTasks()))
		newTaskJson, _ := json.Marshal(&newTask)
		request, _ := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewBuffer(newTaskJson))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		tasksInStore := worktrackerServer.store.GetAllTasks()
		for _, taskInStore := range tasksInStore {
			if taskInStore.Id != newTask.Id {
				assert.False(t, taskInStore.Active)
			} else {
				assert.True(t, taskInStore.Active)
			}
		}
	})

	t.Run("stop current task", func(t *testing.T) {
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(createTasks()))
		request, _ := http.NewRequest(http.MethodPost, "/tasks/stop", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		tasksInStore := worktrackerServer.store.GetAllTasks()
		assertTasksInactive(t, tasksInStore)
	})

	t.Run("start stopped task again", func(t *testing.T) {
		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStoreWithoutTimeIntervals(createTasks()))
		timeNow := time.Now()
		requestData := []byte(fmt.Sprintf(`{
	        "startTime": "%v"
	    }`, timeNow.Format(time.RFC3339)))
		request, _ := http.NewRequest(http.MethodPost, "/tasks/0", bytes.NewBuffer(requestData))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		for _, task := range worktrackerServer.store.GetAllTasks() {
			if task.Id == 0 {
				assert.True(t, task.Active)
			} else {
				assert.False(t, task.Active)
			}
		}
	})

	t.Run("task has correct duration after stop", func(t *testing.T) {
		timeIntervalStart := time.Now()
		timeInterval := &TimeInterval{StartTime: timeIntervalStart, EndTime: nil}
		task := &Task{Id: 0, Title: "Some task", Description: "Some desc", Active: true, Duration: 0}
		tasks := make([]*Task, 0)
		tasks = append(tasks, task)
		timeIntervals := make(map[int][]*TimeInterval)
		timeIntervals[task.Id] = []*TimeInterval{timeInterval}

		timeIntervalEnd := timeIntervalStart.Add(time.Hour)
		patch := monkey.Patch(time.Now, func() time.Time { return timeIntervalEnd })
		defer patch.Unpatch()

		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(tasks, timeIntervals))
		request, _ := http.NewRequest(http.MethodPost, "/tasks/stop", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskFromStore := worktrackerServer.store.GetTaskById(task.Id)
		duration := timeIntervalEnd.Sub(timeIntervalStart)
		assert.Equal(t, duration, taskFromStore.Duration)
	})

	t.Run("task has correct duration after multiple stops", func(t *testing.T) {
		timeStartInterval1 := time.Now()
		timeEndInterval1 := timeStartInterval1.Add(time.Minute * 5)
		timeStartInterval2 := timeEndInterval1.Add(time.Minute * 10)
		timeInterval1 := &TimeInterval{StartTime: timeStartInterval1, EndTime: &timeEndInterval1}
		timeInterval2 := &TimeInterval{StartTime: timeStartInterval2, EndTime: nil}
		task := &Task{Id: 0, Title: "Some task", Description: "Some desc", Active: true, Duration: 0}
		tasks := make([]*Task, 0)
		tasks = append(tasks, task)
		timeIntervals := make(map[int][]*TimeInterval)
		timeIntervals[task.Id] = []*TimeInterval{timeInterval1, timeInterval2}

		timeEndInterval2 := timeStartInterval2.Add(time.Hour)
		timeNowPatch := monkey.Patch(time.Now, func() time.Time { return timeEndInterval2 })
		defer timeNowPatch.Unpatch()

		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(tasks, timeIntervals))
		request, _ := http.NewRequest(http.MethodPost, "/tasks/stop", nil)
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskFromStore := worktrackerServer.store.GetTaskById(task.Id)
		duration := timeEndInterval2.Sub(timeStartInterval2) + timeEndInterval1.Sub(timeStartInterval1)
		assert.Equal(t, duration, taskFromStore.Duration)
	})

	t.Run("task has correct duration after starting new active task", func(t *testing.T) {
		timeIntervalStart := time.Now()
		timeInterval := &TimeInterval{StartTime: timeIntervalStart, EndTime: nil}
		task := &Task{Id: 0, Title: "Some task", Description: "Some desc", Active: true, Duration: 0}
		tasks := make([]*Task, 0)
		tasks = append(tasks, task)
		timeIntervals := make(map[int][]*TimeInterval)
		timeIntervals[task.Id] = []*TimeInterval{timeInterval}

		newTask := &NewTask{
			Title:       "New task",
			Description: "Some desc",
			StartTime:   timeIntervalStart,
			EndTime:     nil,
		}

		timeIntervalEnd := timeIntervalStart.Add(time.Hour)
		patch := monkey.Patch(time.Now, func() time.Time { return timeIntervalEnd })
		defer patch.Unpatch()

		worktrackerServer := NewWorktrackerServer(NewInMemoryWorktrackerStore(tasks, timeIntervals))
		newTaskJson, _ := json.Marshal(&newTask)
		request, _ := http.NewRequest(http.MethodPost, "/tasks/", bytes.NewBuffer(newTaskJson))
		response := httptest.NewRecorder()

		worktrackerServer.ServeHTTP(response, request)

		taskFromStore := worktrackerServer.store.GetTaskById(task.Id)
		duration := timeIntervalEnd.Sub(timeIntervalStart)
		assert.Equal(t, duration, taskFromStore.Duration)
	})
}
