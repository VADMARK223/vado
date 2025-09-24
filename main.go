package main

import (
	"vado/gui"
	"vado/gui/tab/lesson/workers"
)

const showGui = false

func main() {
	if showGui {
		gui.ShowMainApp()
	} else {
		workers.RunWorkersWithContext()
	}
}
