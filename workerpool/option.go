package workerpool

type Option func(p *Pool)

// default number of worker
const noWorker = 1

// default tasks queue capacity
const tasksQueueCap = 100

// use to configuire number of workers in the pool
func WithNoWorkers(i int) Option {
	return func(p *Pool) {
		if i <= 0 {
			return
		}
		p.maxWorkers = i
	}
}

// default capacity of internal job queue of the pool
func WithJobQueueCap(i int) Option {
	return func(p *Pool) {
		if i <= 0 {
			return
		}
		p.tasksqueue = make(chan Task, i)
	}
}
