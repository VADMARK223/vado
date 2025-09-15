package tabs

import (
	gui "vado/gui/common"
	"vado/module/atomic"
	"vado/module/http"
	"vado/module/mutex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func CreateModulesTab() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(gui.CreateBtn("HTTP", func() {
		go http.StartServer()
	}))
	vBox.Add(gui.CreateBtn("RwMutex", mutex.RunRxMutex))
	vBox.Add(gui.CreateBtn("Mutex", mutex.RunMutex))
	vBox.Add(gui.CreateBtn("Atomic", atomic.RunAtomic))
	vBox.Add(gui.CreateBtn("Posts and miners", atomic.RunAtomic, gui.ButtonDisable()))

	return vBox
}
