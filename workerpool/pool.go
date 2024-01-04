package workerpool

import "sync"

type Pool struct {
	// internal queue to submit new task
	tasksqueue chan Task
	// communication channel between worker and the pool
	taskschan chan chan Task
	// list of available worker
	workers []*worker
	// set minimum numbers of workers. Set by WithMinWorkers
	minWorkers int
	// set maximum numbers of workers. Set by WithMaxWorkers
	maxWorkers int
	// wait for worker to finish  their job
	workersWait *sync.WaitGroup
	// to quit all workersk
	quit chan struct{}
}

// kick off all workers
func (p *Pool) Start() {
	for i := range p.workers {
		p.workers[i].start()
	}

	go p.Listen()
}

func (p *Pool) Listen() {
	for {
		select {
		case job := <-p.tasksqueue:
			ch := <-p.taskschan
			ch <- job
		case <-p.quit:
			for i := range p.workers {
				p.workers[i].stop()
			}
			p.workersWait.Wait()
			return
		}
	}
}

// Add task  to the pool for worker to handled
func (p *Pool) AddTask(t Task) {
	p.tasksqueue <- t
}

// if opts are not called , pool will use default value
func New(opts ...Option) *Pool {
	workers := make([]*worker, noWorker, noWorker)
	taskschan := make(chan chan Task)
	workersWait := sync.WaitGroup{}

	for i := range workers {
		workers[i] = newWorker(i+1, taskschan, &workersWait)
	}

	p := &Pool{
		tasksqueue:  make(chan Task, tasksQueueCap),
		taskschan:   taskschan,
		workers:     workers,
		workersWait: &workersWait,
		quit:        make(chan struct{}),
	}
	return p
}
