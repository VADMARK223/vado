package mutex

import (
	"fmt"
	"sync"
	"time"
)

var likes = 0
var rwMtx sync.RWMutex

func setLikes(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10_000; i++ {
		rwMtx.Lock()
		likes++
		rwMtx.Unlock()
	}
}

func getLikes(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 10_000; i++ {
		rwMtx.RLock()
		_ = likes
		rwMtx.RUnlock()
	}
}

func RunRxMutex() {
	fmt.Println("Run RxMutex")
	wg := &sync.WaitGroup{}

	initTime := time.Now()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go setLikes(wg)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go getLikes(wg)
	}

	wg.Wait()

	fmt.Println("Program execution time:", time.Since(initTime))
}
