package database

import (
	"fmt"
	"sync"
	"time"
)

// TaskStatus представляет статус выполнения задачи
type TaskStatus int

const (
	StatusInProgress TaskStatus = iota
	StatusCompleted
	StatusError
)

// Task представляет структуру задачи в базе данных
type Task struct {
	Expression string
	Status     TaskStatus
	Result     string
}

// Database представляет базу данных задач
type Database struct {
	tasks map[string]Task
	mu    sync.Mutex
}

// NewDatabase создает новую базу данных
func NewDatabase() *Database {
	return &Database{
		tasks: make(map[string]Task),
	}
}

// AddTask добавляет новую задачу в базу данных
func (db *Database) AddTask(id string, expression string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Предполагаем, что по умолчанию задача в статусе "в работе"
	db.tasks[id] = Task{
		Expression: expression,
		Status:     StatusInProgress,
	}
}

// GetTask возвращает задачу по её идентификатору
func (db *Database) GetTask(id string) (Task, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	task, exists := db.tasks[id]
	return task, exists
}

// SetTaskResult устанавливает результат выполнения задачи
func (db *Database) SetTaskResult(id string, status TaskStatus, result string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	task, exists := db.tasks[id]
	if exists {
		task.Status = status
		task.Result = result
		db.tasks[id] = task
		return nil
	}
	return fmt.Errorf("task with id %s not found", id)
}

// GetAllTasks возвращает все задачи в базе данных
func (db *Database) GetAllTasks() map[string]Task {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.tasks
}

// GenerateID генерирует уникальный идентификатор задачи
func (db *Database) GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetTaskStatus возвращает статус задачи
func GetTaskStatus(status TaskStatus) string {

	if status == StatusInProgress {
		return "In progress"
	} else if status == StatusCompleted {
		return "Completed"
	} else if status == StatusError {
		return "Error"
	} else {
		return "!unknown status!"
	}
}
