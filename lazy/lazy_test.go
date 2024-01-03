package lazy_test

import (
	"sync"
	"testing"

	"github.com/diontr00/go-helper/lazy"
)

func loadConfig(ch chan bool) func() {
	return func() {
		ch <- true
	}
}

func TestLazy(t *testing.T) {
	n := 5
	ch := make(chan bool)
	loadfn := loadConfig(ch)
	load := lazy.New(loadfn)
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

	for <-ch {
		count++
	}

	if count != 1 {
		t.Errorf("Expect called 1 times , but got %d", count)
	}
}
