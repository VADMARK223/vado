package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Toggle struct {
	widget.BaseWidget
	offBtn    *widget.Button
	onBtn     *widget.Button
	on        *bool
	OnChanged func(*bool)
}

func NewToggleDefault() *Toggle {
	return NewToggle(nil)
}

func NewToggle(initOn *bool) *Toggle {
	onBtn := widget.NewButton("ON", nil)
	offBtn := widget.NewButton("OFF", nil)
	t := &Toggle{offBtn: offBtn, onBtn: onBtn, on: initOn}

	offBtn.OnTapped = func() { t.SetToggle(false) }
	onBtn.OnTapped = func() { t.SetToggle(true) }

	t.ExtendBaseWidget(t)
	t.UpdateUI()

	return t
}

func (t *Toggle) CreateRenderer() fyne.WidgetRenderer {
	box := container.New(&tightHBox{}, t.offBtn, t.onBtn)
	return widget.NewSimpleRenderer(box)
}

func (t *Toggle) SetToggle(value bool) {
	if t.on == nil || *t.on != value {
		t.on = &value
		t.UpdateUI()
		if t.OnChanged != nil {
			t.OnChanged(t.on)
		}
	}
}

func (t *Toggle) UpdateUI() {
	switch {
	case t.on == nil:
		t.offBtn.Importance = widget.MediumImportance
		t.onBtn.Importance = widget.MediumImportance
	case *t.on:
		t.offBtn.Importance = widget.MediumImportance
		t.onBtn.Importance = widget.HighImportance
	default:
		t.offBtn.Importance = widget.HighImportance
		t.onBtn.Importance = widget.MediumImportance
	}
	t.offBtn.Refresh()
	t.onBtn.Refresh()
}

func (t *Toggle) Clear() {
	t.on = nil
	t.UpdateUI()
	if t.OnChanged != nil {
		t.OnChanged(nil)
	}
}

func (t *Toggle) SetOffBtnText(text string) {
	t.offBtn.SetText(text)
}

func (t *Toggle) SetOnBtnText(text string) {
	t.onBtn.SetText(text)
}

type tightHBox struct{}

func (t *tightHBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var width, height float32
	for _, o := range objects {
		size := o.MinSize()
		width += size.Width
		if size.Height > height {
			height = size.Height
		}
	}
	return fyne.NewSize(width, height)
}

func (t *tightHBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	x := float32(0)
	for _, o := range objects {
		oMin := o.MinSize()
		o.Resize(oMin)
		o.Move(fyne.NewPos(x, (size.Height-oMin.Height)/2))
		x += oMin.Width
	}
}
