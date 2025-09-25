package waitGroup

import (
	"fmt"
	"sync"
)

func RunWaitGroup() {
	fmt.Println("Running WaitGroup")
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		i := i
		go func() {
			defer wg.Done()
			fmt.Println("i:", i)
		}()
	}

	wg.Wait()
}
