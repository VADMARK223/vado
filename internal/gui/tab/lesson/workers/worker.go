package workers

import (
	"context"
	"fmt"
	"time"
	"vado/pkg/util"
)

type Worker struct {
	name  string
	value int
}

func (w *Worker) Work(channel chan<- int) {
	workTime := util.RndIntn(4) + 1
	fmt.Printf("  %s start work (%d sec)...\n", w.name, workTime)
	time.Sleep(time.Duration(workTime) * time.Second)
	fmt.Printf("  %s write in channel.\n", w.name)
	channel <- w.value
	fmt.Printf("  %s stop work.\n", w.name)
}

func (w *Worker) WorkWithContext(ctx context.Context, ch chan<- int) {
	workTime := util.RndIntn(4) + 1
	fmt.Printf("  %s start work (%d sec)...\n", w.name, workTime)

	select {
	case <-time.After(time.Duration(workTime) * time.Second): // Рабочий доработал. Сработает, когда timeWork истечет
		fmt.Printf("  %s finished work, writing in channel.\n", w.name)
		// Если сделать без select, и канал никто не читает, тут будет блокировка.
		select {
		case ch <- w.value:
			fmt.Printf("  %s wrote result.\n", w.name)
		case <-ctx.Done():
			// Отменили пока ждал запись
			fmt.Printf("  %s canceled (timeout).\n", w.name)
		}
	case <-ctx.Done(): // Рабочий не успел доработать
		fmt.Printf("  %s canceled before finishing.\n", w.name)
	}
}
