package tasks

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	c "vado/constant"
	m "vado/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateTasksGui() fyne.CanvasObject {
	title := canvas.NewText("Tasks list", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	data, err := os.ReadFile(c.TaskFileName)

	if err != nil {
		log.Fatal(err)
	}

	var list m.TaskList
	if err = json.Unmarshal(data, &list); err != nil {
		fmt.Printf("Error: '%s'\n", err)
	}
	vBox := container.NewVBox()
	for _, task := range list.Tasks {
		vBox.Add(widget.NewLabel(task.Name))
	}

	scroll := container.NewVScroll(vBox)
	content := container.NewBorder(title, nil, nil, nil, scroll)

	return content
}
