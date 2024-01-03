package lazy

import "sync"

type lazy struct {
	load func()
	once sync.Once
}

// Load the wrapped load function, its guarantee to run only one
func (l *lazy) Load() {
	l.once.Do(l.load)
}

// Create new lazy by wrap around load function
func New(load func()) *lazy {
	return &lazy{
		load: load,
	}
}
