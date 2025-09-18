package atomic

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var number atomic.Int64

func increase(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		number.Add(1)
	}
}

func RunAtomic() {
	wg := &sync.WaitGroup{}
	fmt.Println("Start atomic")
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Add(1)
	go increase(wg)
	wg.Wait()
	fmt.Println("Result number", number.Load())
	number.Store(0)
	fmt.Println("End atomic")
}
