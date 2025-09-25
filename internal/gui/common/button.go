package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ButtonOption func(*widget.Button)

func ButtonDisable() ButtonOption {
	return func(btn *widget.Button) {
		btn.Disable()
	}
}

// CreateBtn создает кнопку
func CreateBtn(label string, icon fyne.Resource, tapped func(), opts ...ButtonOption) *widget.Button {
	btn := widget.NewButtonWithIcon(label, icon, func() {})
	btn.OnTapped = tapped

	for _, opt := range opts {
		opt(btn)
	}
	return btn
}
