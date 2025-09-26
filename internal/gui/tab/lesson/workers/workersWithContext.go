package workers

import (
	"context"
	"fmt"
	"sync"
)

func RunWorkersWithContext() {
	ctx, cancel := context.WithTimeout(context.Background(), workerTimeout)
	defer cancel()

	ch := make(chan int, numWorkers)
	var wg sync.WaitGroup

	workers := []*Worker{
		{name: "Worker #1", value: 1},
		{name: "Worker #2", value: 2},
		{name: "Worker #3", value: 3},
	}

	for _, w := range workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			w.WorkWithContext(ctx, ch)
		}(w)
	}

	// В отдельной горутине, чтобы цикл после читал канал
	go func() {
		wg.Wait() // Блокирует горутину, закрытие канала не будет, пока все workers не завершаться.
		close(ch) // Когда все workers закончили → закрываем канал (это сообщит v := range ch что можно уже не читать)
	}()

	result := 0
	for v := range ch {
		result += v
	}

	fmt.Println("Result:", result)
	fmt.Println("Main exit.")
}
