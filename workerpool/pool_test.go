package workerpool_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/diontr00/go-helper/workerpool"
)

type testTask struct {
	expectError  bool
	errHandled   bool
	valueHandled bool
	wg           *sync.WaitGroup
}

func (t *testTask) Execute() error {
	if t.expectError {
		return errors.New("Expect return error")
	}
	t.valueHandled = true
	t.wg.Done()
	return nil
}
func (t *testTask) OnErrorHandle(err error) {
	t.errHandled = true
	t.wg.Done()
}

func newTestTask(expectErr bool, wg *sync.WaitGroup) *testTask {
	return &testTask{expectError: expectErr, wg: wg}
}

func TestCreatePoolWithDefault(t *testing.T) {
	p := workerpool.New()
	wg := sync.WaitGroup{}

	t1 := newTestTask(false, &wg)
	t2 := newTestTask(true, &wg)

	time.Sleep(time.Second)

	p.Start()
	wg.Add(2)
	p.AddTask(t1)
	p.AddTask(t2)
	wg.Wait()

	if !t1.valueHandled {
		t.Error("Expect handled value for task 1")
	}

	if !t2.errHandled {
		t.Error("Expect handled error for task 2")
	}
}
