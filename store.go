package main

type WorktrackerStore interface{
	GetAllTasks() []*Task
	GetTaskById(taskId int) *Task
}

type InMemoryWorktrackerStore struct {
	tasks map[int]*Task
}

func NewInMemoryWorktrackerStore(tasks []*Task) *InMemoryWorktrackerStore {
	return &InMemoryWorktrackerStore{createTasksMap(tasks)}
}

func (s *InMemoryWorktrackerStore) GetAllTasks() []*Task {
	return getTasksList(s.tasks)
}

func (s *InMemoryWorktrackerStore) GetTaskById(taskId int) *Task {
	return s.tasks[taskId]
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
	for  _, value := range tasks {
		tasksList = append(tasksList, value)
	}
	return tasksList;
}