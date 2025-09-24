package workers

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Worker struct {
	name  string
	value int
}

func (w *Worker) Work(channel chan<- int) {
	workTime := rand.IntN(4) + 1
	fmt.Printf("  %s start work (%d sec)...\n", w.name, workTime)
	time.Sleep(time.Duration(workTime) * time.Second)
	fmt.Printf("  %s write in channel.\n", w.name)
	channel <- w.value
	fmt.Printf("  %s stop work.\n", w.name)
}
