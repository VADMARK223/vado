package main

import (
	"vado/gui"
	"vado/gui/tab/lesson/waitGroup"
)

const showGui = false

func main() {
	if showGui {
		gui.ShowMainApp()
	} else {
		waitGroup.RunWaitGroup()
	}
}
