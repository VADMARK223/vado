package main

import (
	"vado/gui"
	"vado/gui/tab/lesson/database"
)

const showGui = true

func main() {
	if showGui {
		gui.ShowMainApp()
	} else {
		database.RunDatabase()
	}
}
