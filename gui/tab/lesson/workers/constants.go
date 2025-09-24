package workers

import "time"

const (
	numWorkers    = 3               // Кол-во рабочих горутин
	workerTimeout = 5 * time.Second // Если какой-то worker будет долго тупить с записью в канал, то выводить предупреждение и ждать дальше.
)
