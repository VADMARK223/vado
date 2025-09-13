package modules

import (
	"vado/module/atomic"
	"vado/module/mutex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createButton(label string, tapped func()) *widget.Button {
	btn := widget.NewButton(label, func() {})
	btn.OnTapped = tapped
	return btn
}

func CreateModulesGui() fyne.CanvasObject {
	vBox := container.NewVBox()
	vBox.Add(createButton("RwMutex", mutex.RunRxMutex))
	vBox.Add(createButton("Mutex", mutex.RunMutex))
	vBox.Add(createButton("Atomic", atomic.RunAtomic))

	return vBox
}
