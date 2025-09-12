package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	c "vado/constant"
	m "vado/model"
	"vado/module/atomic"
	"vado/module/mutex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	data, err := os.ReadFile(c.TaskFileName)

	if err != nil {
		log.Fatal(err)
	}

	var list m.TaskList
	if err = json.Unmarshal(data, &list); err != nil {
		fmt.Printf("Error: '%s'\n", err)
	}

	a := app.New()
	w := a.NewWindow("Vado")

	btnExit := widget.NewButton("Exit", func() {
		w.Close()
	})

	title := canvas.NewText("Tasks list", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter
	content := container.NewVBox(title)

	for _, task := range list.Tasks {
		content.Add(widget.NewLabel(task.Name))
	}
	content.Add(widget.NewButton("Atomic", func() {
		atomic.RunAtomic()
	}))

	content.Add(widget.NewButton("Mutex", func() {
		mutex.RunMutex()
	}))
	content.Add(btnExit)

	w.SetContent(content)

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			w.Close()
		}
	})
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
