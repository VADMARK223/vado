package mutex

import (
	"fmt"
	"sync"
)

var slice []int
var mtx sync.Mutex

func increase(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		mtx.Lock()
		slice = append(slice, 1)
		mtx.Unlock()
	}
}

func RunMutex() {
	wg := &sync.WaitGroup{}
	fmt.Println("Start mutex")
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
	fmt.Println("Result len:", len(slice))
	slice = nil
	fmt.Println("End mutex")
}
