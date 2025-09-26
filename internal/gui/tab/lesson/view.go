package lesson

import (
	gui "vado/internal/gui/common"
	"vado/internal/gui/tab/lesson/atomic"
	"vado/internal/gui/tab/lesson/inMemoryCache"
	mutex2 "vado/internal/gui/tab/lesson/mutex"
	"vado/internal/gui/tab/lesson/points"
	"vado/internal/gui/tab/lesson/sliceArray"
	"vado/internal/gui/tab/lesson/waitGroup"
	"vado/internal/gui/tab/lesson/workers"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateView() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(gui.NewBtn("3 workers", nil, func() { run(workers.RunWorkers) }))
	vBox.Add(gui.NewBtn("Slice and array", nil, sliceArray.RunSliceArray))
	vBox.Add(gui.NewBtn("Wait group", nil, waitGroup.RunWaitGroup))
	vBox.Add(gui.NewBtn("Pointers", nil, points.RunPointers))
	vBox.Add(gui.NewBtn("In-memory cache", nil, inMemoryCache.RunInMemoryCache))
	vBox.Add(gui.NewBtn("RwMutex", nil, mutex2.RunRxMutex))
	vBox.Add(gui.NewBtn("Mutex", nil, mutex2.RunMutex))
	vBox.Add(gui.NewBtn("Atomic", nil, atomic.RunAtomic))
	vBox.Add(gui.NewBtn("Posts and miners", nil, atomic.RunAtomic, gui.ButtonDisable()))

	return vBox
}

func run(task func()) {
	go task()
}
