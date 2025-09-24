package workers

import (
	"fmt"
	"sync"
	"time"
)

func RunWorkers() {
	workers := []*Worker{
		{name: "Worker #1", value: 1},
		{name: "Worker #2", value: 2},
		{name: "Worker #3", value: 3},
	}

	var wg sync.WaitGroup // Не обязательное держать в виде указателя (&sync.WaitGroup{} тоже работает)
	wg.Add(numWorkers)

	// Буферизованный канал на количество workers — чтобы они могли закончить
	ch := make(chan int, numWorkers)

	for _, worker := range workers {
		go runWork(worker, &wg, ch)
	}

	// Общий таймаут на сбор всех результатов
	timeout := time.After(workerTimeout)
	got := 0
	done := false
	result := 0

	for !done || got < numWorkers {
		fmt.Printf("Waiting value #%d...\n", got+1)
		select {
		case v := <-ch:
			result += v
			got++
		case <-timeout:
			fmt.Println("Timeout! Stop waiting for more results")
			done = true
		}
	}

	// Ждём всех workers
	wg.Wait()
	close(ch) // Теперь можно безопасно читать из канала до конца

	// Читаем всё, что осталось в буфере
	for v := range ch {
		fmt.Println("Late value from buffer:", v)
		result += v
	}

	fmt.Println("Result:", result)
}

func runWork(w *Worker, wg *sync.WaitGroup, ch chan<- int) {
	defer wg.Done()
	w.Work(ch)
}
