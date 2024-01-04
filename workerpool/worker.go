package workerpool

import (
	"log"
	"runtime"
	"sync"
)

type worker struct {
	// to keep track the worker
	id int
	// decrease wait to  signal worker exit
	wait *sync.WaitGroup
	// Communication channel between pool and worker
	taskschan chan chan Task
	// worker own tasks chan, which will be push given to the
	// tasksChan  to receive task
	task chan Task
	//  Signal worker stop
	quit chan struct{}
}

// return new worker
func newWorker(id int, taskschan chan chan Task, wait *sync.WaitGroup) *worker {
	return &worker{
		id:        id,
		taskschan: taskschan,
		wait:      wait,
		quit:      make(chan struct{}),
		task:      make(chan Task),
	}
}

// process the start , if panic print stack trace , else error using Errhandler
func (w *worker) process(task Task) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("panic running task: %v\n%s\n", r, buf)
		}
	}()

	err := task.Execute()
	if err != nil {
		task.OnErrorHandle(err)
	}
}

// stop worker
func (w *worker) stop() {
	close(w.quit)
}

// start the worker, worker will send back it own channel to receive for task
func (w *worker) start() {
	go func() {
		defer w.wait.Done()
		w.wait.Add(1)
		for {
			w.taskschan <- w.task
			select {
			// got task
			case task := <-w.task:
				w.process(task)
			case <-w.quit:
				w.wait.Done()
				return
			}
		}
	}()
}
