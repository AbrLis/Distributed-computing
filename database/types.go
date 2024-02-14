package database

import (
	"sync"
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

// TaskStatus представляет статус выполнения задачи
type TaskStatus int
