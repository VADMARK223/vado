package tab

import (
	gui "vado/gui/common"
	"vado/module/atomic"
	"vado/module/mutex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateModulesTab() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(gui.CreateBtn("RwMutex", nil, mutex.RunRxMutex))
	vBox.Add(gui.CreateBtn("Mutex", nil, mutex.RunMutex))
	vBox.Add(gui.CreateBtn("Atomic", nil, atomic.RunAtomic))
	vBox.Add(gui.CreateBtn("Posts and miners", nil, atomic.RunAtomic, gui.ButtonDisable()))

	return vBox
}
