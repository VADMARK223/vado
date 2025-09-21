package main

import (
	"vado/gui"
	"vado/gui/tab/lesson/inMemoryCache"
)

const showGui = false

func main() {
	if showGui {
		gui.ShowMainApp()
	} else {
		inMemoryCache.RunInMemoryCache()
	}
}
