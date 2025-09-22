package lesson

import (
	gui "vado/gui/common"
	"vado/gui/tab/lesson/atomic"
	"vado/gui/tab/lesson/inMemoryCache"
	"vado/gui/tab/lesson/mutex"
	"vado/gui/tab/lesson/points"
	"vado/gui/tab/lesson/sliceArray"
	"vado/gui/tab/lesson/waitGroup"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateView() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(gui.CreateBtn("Slice and array", nil, sliceArray.RunSliceArray))
	vBox.Add(gui.CreateBtn("Wait group", nil, waitGroup.RunWaitGroup))
	vBox.Add(gui.CreateBtn("Pointers", nil, points.RunPointers))
	vBox.Add(gui.CreateBtn("In-memory cache", nil, inMemoryCache.RunInMemoryCache))
	vBox.Add(gui.CreateBtn("RwMutex", nil, mutex.RunRxMutex))
	vBox.Add(gui.CreateBtn("Mutex", nil, mutex.RunMutex))
	vBox.Add(gui.CreateBtn("Atomic", nil, atomic.RunAtomic))
	vBox.Add(gui.CreateBtn("Posts and miners", nil, atomic.RunAtomic, gui.ButtonDisable()))

	return vBox
}
