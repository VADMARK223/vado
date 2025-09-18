package lesson

import (
	gui "vado/gui/common"
	"vado/gui/tab/lesson/atomic"
	"vado/gui/tab/lesson/database"
	mutex2 "vado/gui/tab/lesson/mutex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateView() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(gui.CreateBtn("Database", nil, database.RunDatabase))
	vBox.Add(gui.CreateBtn("RwMutex", nil, mutex2.RunRxMutex))
	vBox.Add(gui.CreateBtn("Mutex", nil, mutex2.RunMutex))
	vBox.Add(gui.CreateBtn("Atomic", nil, atomic.RunAtomic))
	vBox.Add(gui.CreateBtn("Posts and miners", nil, atomic.RunAtomic, gui.ButtonDisable()))

	return vBox
}
