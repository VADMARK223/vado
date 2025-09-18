package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "fyne.io/fyne/v2/widget"
)

// CreateLazyTabItem создает вкладку с ленивой загрузкой содержимого.
//
// Созданная вкладка изначально отображает индикатор прогресса.
// Настоящее содержимое создается только при первой активации вкладки.
//
// Параметры:
//
//	text — заголовок вкладки.
//	factory — функция, которая возвращает объект fyne.CanvasObject для содержимого вкладки.
//	factories — карта, в которую сохраняется связь между вкладкой и функцией factory; используется для ленивой инициализации.
//
// Возвращает:
//
//	Указатель на созданную вкладку *container.TabItem.
func CreateLazyTabItem(text string,
	factory func() fyne.CanvasObject,
	factories map[*container.TabItem]func() fyne.CanvasObject) *container.TabItem {
	placeholder := widget.NewProgressBarInfinite()
	tab := container.NewTabItem(text, placeholder)
	factories[tab] = factory
	return tab
}
