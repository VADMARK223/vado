package tabs

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	c "vado/constant"
	"vado/gui/common"
	m "vado/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CreateTasksTab(win fyne.Window) fyne.CanvasObject {
	title := canvas.NewText("Задания", color.White)
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

	for i := range list.Tasks {
		task := list.Tasks[i]

		taskLabel := widget.NewLabel(fmt.Sprintf("%d %s", task.Id, task.Name))
		hBox := container.NewHBox(taskLabel, layout.NewSpacer())
		taskDelBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
			if common.ShowDeleteConfirm {
				dialog.ShowConfirm("Удаление задания", "Выдействительно хотите удалить задание?", func(b bool) {
					if b {
						DeleteTask(&list, task.Id, hBox, vBox)
					}
				}, win)
			} else {
				DeleteTask(&list, task.Id, hBox, vBox)
			}
		})
		hBox.Add(taskDelBtn)
		vBox.Add(hBox)
	}

	scroll := container.NewVScroll(vBox)
	content := container.NewBorder(title, nil, nil, nil, scroll)

	return content
}

func DeleteTask(list *m.TaskList, taskId int, hBox *fyne.Container, vBox *fyne.Container) {
	var newTasks []m.Task
	for _, t := range list.Tasks {
		if t.Id != taskId {
			newTasks = append(newTasks, t)
		}
	}
	list.Tasks = newTasks

	// Сохраняем в json
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		fmt.Println("Error marshal JSON:", err)
		return
	}
	/*
		0 → это префикс, который говорит Go, что число в восьмеричной системе.
		6 → владелец может читать и писать (4+2).
		4 → группа может только читать.
		4 → остальные пользователи могут только читать.
		os.ModePerm - это 0777, то есть маскимально открытые права.
	*/
	if err := os.WriteFile(c.TaskFileName, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	vBox.Remove(hBox)
	vBox.Refresh()
}
