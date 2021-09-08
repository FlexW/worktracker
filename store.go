package main

import "time"

type WorktrackerStore interface {
	GetAllTasks() []*Task
	GetTaskById(taskId int) *Task
	InsertTask(task *Task) int
	UpdateTask(task *Task)
	AddTimeIntervalToTask(taskId int, timeInterval *TimeInterval)
	GetTimeIntervalsByTaskId(taskId int) []*TimeInterval
	SetTaskInactive(taskId int)
	SetTaskActive(taskId int)
}

type InMemoryWorktrackerStore struct {
	tasks         map[int]*Task
	timeIntervals map[int][]*TimeInterval
}

func NewInMemoryWorktrackerStore(tasks []*Task, timeIntervals map[int][]*TimeInterval) *InMemoryWorktrackerStore {
	return &InMemoryWorktrackerStore{tasks: createTasksMap(tasks), timeIntervals: timeIntervals}
}

func NewInMemoryWorktrackerStoreWithoutTimeIntervals(tasks []*Task) *InMemoryWorktrackerStore {
	return &InMemoryWorktrackerStore{tasks: createTasksMap(tasks), timeIntervals: make(map[int][]*TimeInterval)}
}

func (s *InMemoryWorktrackerStore) GetAllTasks() []*Task {
	return getTasksList(s.tasks)
}

func (s *InMemoryWorktrackerStore) GetTaskById(taskId int) *Task {
	return s.tasks[taskId]
}

func (s *InMemoryWorktrackerStore) InsertTask(task *Task) int {
	id := len(s.tasks)
	s.tasks[id] = task
	return id
}

func (s *InMemoryWorktrackerStore) UpdateTask(task *Task) {
	s.tasks[task.Id] = task
}

func calculateDuration(timeIntervals []*TimeInterval) time.Duration {
	var entireDuration time.Duration
	for _, timeInterval := range timeIntervals {
		if timeInterval.EndTime != nil {
			entireDuration += timeInterval.EndTime.Sub(timeInterval.StartTime)
		} else {
			entireDuration += time.Now().Sub(timeInterval.StartTime)
		}
	}
	return entireDuration

}

func (s *InMemoryWorktrackerStore) updateTaskDuration(taskId int) {
	taskFromStore := s.tasks[taskId]
	timeIntervals := s.timeIntervals[taskId]
	taskFromStore.Duration = calculateDuration(timeIntervals)
}

func (s *InMemoryWorktrackerStore) AddTimeIntervalToTask(taskId int, timeInterval *TimeInterval) {
	timeIntervals := s.timeIntervals[taskId]
	timeIntervals = append(timeIntervals, timeInterval)
	s.updateTaskDuration(taskId)
}

func (s *InMemoryWorktrackerStore) GetTimeIntervalsByTaskId(taskId int) []*TimeInterval {
	return s.timeIntervals[taskId]
}

func (s *InMemoryWorktrackerStore) SetTaskInactive(taskId int) {
	s.tasks[taskId].Active = false

	timeIntervals := s.timeIntervals[taskId]
	for i := range timeIntervals {
		if timeIntervals[i].EndTime == nil {
			timeNow := time.Now()
			timeIntervals[i].EndTime = &timeNow
		}
	}
	s.updateTaskDuration(taskId)
}

func (s *InMemoryWorktrackerStore) SetTaskActive(taskId int) {
	s.tasks[taskId].Active = true
}

func createTasksMap(tasks []*Task) map[int]*Task {
	taskMap := make(map[int]*Task)
	for _, task := range tasks {
		taskMap[task.Id] = task
	}
	return taskMap
}

func getTasksList(tasks map[int]*Task) []*Task {
	tasksList := make([]*Task, 0, len(tasks))
	for _, value := range tasks {
		tasksList = append(tasksList, value)
	}
	return tasksList
}
