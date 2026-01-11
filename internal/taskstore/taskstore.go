package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStoreInterface interface {
	CreateTask(text string, tags []string, due time.Time) int
	GetTask(id int) (Task, error)
	DeleteTask(id int) error
	DeleteAllTask() error
	GetAllTasks() []Task
	GetTaskByTag(tag string) []Task
	GetTaskByDueDate(year int, month time.Month, day int) []Task
}

// TaskStore is a simple in-memory database of tasks;
// TaskStore methods are safe to call concurrently.
type TaskStore struct {
	sync.Mutex

	tasks  map[int]Task
	nextId int
}

func (ts *TaskStore) CreateTask(text string, tags []string) int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  time.Now(),
	}

	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++

	return task.Id
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	if task, ok := ts.tasks[id]; ok {
		return task, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id %d not found", id)
	}

	delete(ts.tasks, id)

	return nil
}

func (ts *TaskStore) DeleteAllTask() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	allTask := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		allTask = append(allTask, task)
	}

	return allTask
}

func (ts *TaskStore) GetTaskByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		for _, currTag := range task.Tags {
			if currTag == tag {
				tasks = append(tasks, task)
			}
		}
	}

	return tasks
}

func (ts *TaskStore) GetTaskByDueDate(year int, month time.Month, day int) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		y, m, d := task.Due.Date()

		if y == year && m == month && d == day {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]Task),
		nextId: 0,
	}
}
