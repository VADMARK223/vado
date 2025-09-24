package main

import (
	"vado/gui"
	"vado/gui/tab/lesson/workers"
)

const showGui = true

func main() {
	if showGui {
		gui.ShowMainApp()
	} else {
		workers.RunWorkers()
	}
}
