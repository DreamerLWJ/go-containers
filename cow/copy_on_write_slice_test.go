package cow

import (
	"sync"
	"testing"
)

func TestCopyOnWriteSlice_Append(t *testing.T) {
	test2 := NewCopyOnWriteSlice(0)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			for j := i * 10; j < (i+1)*10; j++ {
				test2.Append(j)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

}
