package settings

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// CreateView представление вкладки настроек приложения
//
// При локальной отладке сохраняет в /home/vadmark/.var/app/com.jetbrains.GoLand/config/fyne/io.vado
func CreateView() fyne.CanvasObject {
	label := widget.NewLabel("В разработке...")
	prefs := fyne.CurrentApp().Preferences()
	log.Println(prefs.String("tasks_file_path"))
	prefs.SetString("tasks_file_path", "./data/tasks.json")

	temp := prefs.String("tasks_file_path")
	log.Println(temp)
	return label
}
