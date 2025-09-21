package lesson

import (
	gui "vado/gui/common"
	"vado/gui/tab/lesson/atomic"
	"vado/gui/tab/lesson/inMemoryCache"
	"vado/gui/tab/lesson/mutex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateView() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(gui.CreateBtn("In-memory cache", nil, inMemoryCache.RunInMemoryCache))
	vBox.Add(gui.CreateBtn("RwMutex", nil, mutex.RunRxMutex))
	vBox.Add(gui.CreateBtn("Mutex", nil, mutex.RunMutex))
	vBox.Add(gui.CreateBtn("Atomic", nil, atomic.RunAtomic))
	vBox.Add(gui.CreateBtn("Posts and miners", nil, atomic.RunAtomic, gui.ButtonDisable()))

	return vBox
}
