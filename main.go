package main

import (
	"vado/gui"
	"vado/gui/tab/http"
)

const showGui = true

func main() {
	if showGui {
		gui.ShowMainApp()
	} else {
		http.StartServer()
	}
}
