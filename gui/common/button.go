package common

import "fyne.io/fyne/v2/widget"

type ButtonOption func(*widget.Button)

func ButtonDisable() ButtonOption {
	return func(btn *widget.Button) {
		btn.Disable()
	}
}

func CreateBtn(label string, tapped func(), opts ...ButtonOption) *widget.Button {
	btn := widget.NewButton(label, func() {})
	btn.OnTapped = tapped

	for _, opt := range opts {
		opt(btn)
	}
	return btn
}
