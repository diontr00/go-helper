package lazy

import "sync"

// Represent lazy load task
type Task interface {
	// Load function that load data
	Load()
}

type lazy struct {
	task Task
	once sync.Once
}

// Load the wrapped load function, its guarantee to run only one
func (l *lazy) Load() {
	l.once.Do(l.task.Load)
}

// Create new lazy by wrap around load function
func New(task Task) *lazy {
	return &lazy{
		task: task,
	}
}
