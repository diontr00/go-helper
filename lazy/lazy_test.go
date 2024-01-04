package lazy_test

import (
	"sync"
	"testing"

	"github.com/diontr00/go-helper/lazy"
)

type testTask struct {
	ch chan struct{}
}

func (t *testTask) Load() {
	t.ch <- struct{}{}
}

func newTask(ch chan struct{}) *testTask {
	return &testTask{
		ch: ch,
	}
}

func TestLazy(t *testing.T) {
	n := 5
	ch := make(chan struct{})

	task := newTask(ch)

	load := lazy.New(task)
	var wg sync.WaitGroup

	count := 0
	go func() {
		wg.Wait()
		close(ch)
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			load.Load()
		}()
	}

	for {
		_, ok := <-ch
		if !ok {
			break
		}
		count++
	}

	if count != 1 {
		t.Errorf("Expect called 1 times , but got %d", count)
	}
}
